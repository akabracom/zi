import { useEffect, useMemo, useState } from 'react'
import {
  getDraftRequest,
  removeServiceFromRequest,
  submitRequest,
} from '../api/http'
import type { Request } from '../api/types'

// Флаг для проверки mock режима
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true' || false

// Локальное состояние для редактирования дат периода
interface EditableItem {
  serviceId: number
  startDate: number // День начала (1-Days)
  endDate: number   // День окончания (1-Days)
}

export function DraftPage() {
  const [draft, setDraft] = useState<Request | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [info, setInfo] = useState<string | null>(null)
  const [commonAmount, setCommonAmount] = useState(100000) // Общая сумма вклада
  const [editableItems, setEditableItems] = useState<Map<number, EditableItem>>(new Map())

  async function loadDraft() {
    setLoading(true)
    setError(null)
    try {
      const data = await getDraftRequest()
      setDraft(data)
      
      // Инициализируем редактируемые значения из данных
      // По умолчанию: startDate = 1, endDate = Days
      const itemsMap = new Map<number, EditableItem>()
      data.services.forEach((item) => {
        const days = extractDays(item.service.description)
        itemsMap.set(item.service_id, {
          serviceId: item.service_id,
          startDate: 1, // По умолчанию с 1-го числа
          endDate: days, // По умолчанию до последнего дня месяца
        })
      })
      setEditableItems(itemsMap)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Не удалось загрузить расчет',
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadDraft()
  }, [])

  // Извлекаем количество дней из description (формат: "Дней: 31 | ...")
  function extractDays(description?: string): number {
    if (!description) return 30
    const match = description.match(/Дней:\s*(\d+)/i)
    return match ? parseInt(match[1], 10) : 30
  }

  // Рассчитываем прибыль для одного месяца: сумма * (ставка/100) * (дней/365)
  function calculateProfit(amount: number, rate: number, days: number): number {
    return amount * (rate / 100) * (days / 365)
  }

  // Получаем редактируемые значения для услуги
  function getEditableItem(serviceId: number): EditableItem {
    const item = editableItems.get(serviceId)
    if (item) return item
    
    // Если нет в состоянии, создаем из draft
    const draftItem = draft?.services.find(s => s.service_id === serviceId)
    if (draftItem) {
      const days = extractDays(draftItem.service.description)
      return {
        serviceId,
        startDate: 1,
        endDate: days,
      }
    }
    
    return { serviceId, startDate: 1, endDate: 30 }
  }

  // Обновление даты начала
  function handleStartDateChange(serviceId: number, newStartDate: number, maxDays: number) {
    if (newStartDate < 1) newStartDate = 1
    if (newStartDate > maxDays) newStartDate = maxDays
    
    const item = getEditableItem(serviceId)
    const updated = new Map(editableItems)
    
    // Если новая дата начала больше даты окончания, обновляем дату окончания
    let endDate = item.endDate
    if (newStartDate > endDate) {
      endDate = newStartDate
    }
    
    updated.set(serviceId, { ...item, startDate: newStartDate, endDate })
    setEditableItems(updated)
  }

  // Обновление даты окончания
  function handleEndDateChange(serviceId: number, newEndDate: number, maxDays: number) {
    if (newEndDate < 1) newEndDate = 1
    if (newEndDate > maxDays) newEndDate = maxDays
    
    const item = getEditableItem(serviceId)
    const updated = new Map(editableItems)
    
    // Если новая дата окончания меньше даты начала, обновляем дату начала
    let startDate = item.startDate
    if (newEndDate < startDate) {
      startDate = newEndDate
    }
    
    updated.set(serviceId, { ...item, startDate, endDate: newEndDate })
    setEditableItems(updated)
  }

  // Вычисляем количество дней между датами
  function calculateDaysBetween(start: number, end: number): number {
    if (end < start) return 0
    return end - start + 1
  }

  const totals = useMemo(() => {
    if (!draft || draft.services.length === 0) {
      return { totalAmount: 0, totalProfit: 0, averageRate: 0 }
    }

    let totalAmount = 0
    let totalProfit = 0
    let totalQuantity = 0
    let weightedSum = 0

    draft.services.forEach((item) => {
      const editable = getEditableItem(item.service_id)
      const periodDays = calculateDaysBetween(editable.startDate, editable.endDate)
      const itemAmount = commonAmount * item.quantity
      const profit = calculateProfit(commonAmount, item.service.price, periodDays) * item.quantity

      totalAmount += itemAmount
      totalProfit += profit
      totalQuantity += item.quantity
      weightedSum += item.quantity * item.service.price
    })

    const averageRate = totalQuantity > 0 ? weightedSum / totalQuantity : 0

    return { totalAmount, totalProfit, averageRate }
  }, [draft, commonAmount, editableItems])

  async function handleRemove(requestId: number, serviceId: number) {
    setError(null)
    if (!confirm('Удалить этот вклад из заявки?')) return
    try {
      await removeServiceFromRequest(requestId, serviceId)
      await loadDraft()
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Не удалось удалить услугу',
      )
    }
  }

  async function handleSubmit() {
    if (!draft) return
    setError(null)
    setInfo(null)
    try {
      await submitRequest(draft.id)
      setInfo('Расчет отправлен на обработку')
      // В mock режиме обновляем статус локально
      if (USE_MOCK && draft) {
        setDraft({ ...draft, status: 'formed' })
      } else {
        await loadDraft()
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Не удалось отправить расчет',
      )
    }
  }

  return (
    <div className="page" style={{ maxWidth: '1000px', margin: '0 auto', padding: '20px' }}>
      <h2 style={{ color: '#009BFF', fontSize: '26px', textAlign: 'center', marginBottom: '25px', fontWeight: '700' }}>
        Выбранные вклады
      </h2>
      
      {loading && <p>Загрузка...</p>}
      {error && <div className="form-error">{error}</div>}
      {info && <div className="form-success">{info}</div>}

      {!loading && !draft && (
        <p className="muted" style={{ textAlign: 'center', padding: '40px' }}>
          Корзина пуста. Добавьте месяцы на странице &quot;Банковские месяцы&quot;.
        </p>
      )}

      {draft && (
        <>
          {/* Поле для общей суммы вклада */}
          <div style={{
            marginBottom: '20px',
            padding: '16px 20px',
            backgroundColor: '#ffffff',
            borderRadius: '10px',
            border: '1px solid #d6ecff',
            boxShadow: '0 2px 6px rgba(0, 155, 255, 0.06)'
          }}>
            <label htmlFor="common-sum" style={{
              fontSize: '14px',
              fontWeight: '600',
              color: '#47607d',
              display: 'block',
              marginBottom: '8px'
            }}>
              Сумма вклада (для всех вкладов):
            </label>
            <input
              type="number"
              id="common-sum"
              min="0"
              step="1000"
              value={commonAmount}
              onChange={(e) => setCommonAmount(Number(e.target.value))}
              style={{
                padding: '6px 10px',
                borderRadius: '6px',
                border: '1px solid #c4d9f5',
                fontSize: '14px',
                maxWidth: '180px',
                backgroundColor: '#ffffff',
                color: '#333333'
              }}
            />
          </div>

          {/* Карточки вкладов */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', marginBottom: '30px' }}>
            {draft.services.map((item) => {
              const editable = getEditableItem(item.service_id)
              const maxDays = extractDays(item.service.description)
              const periodDays = calculateDaysBetween(editable.startDate, editable.endDate)
              const profit = calculateProfit(commonAmount, item.service.price, periodDays) * item.quantity
              
              return (
                <div
                  key={item.service_id}
                  style={{
                    display: 'flex',
                    alignItems: 'flex-start',
                    gap: '20px',
                    padding: '20px',
                    backgroundColor: '#fff',
                    border: '1px solid #e8e8e8',
                    borderRadius: '10px',
                    boxShadow: '0 2px 6px rgba(0,0,0,0.05)'
                  }}
                >
                  {/* Изображение */}
                  {item.service.image_url && (
                    <div>
                      {USE_MOCK ? (
                        <div style={{
                          width: '80px',
                          height: '80px',
                          backgroundColor: '#e0e0e0',
                          borderRadius: '8px',
                          border: '1px solid #ddd',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          color: '#666',
                          fontSize: '12px',
                          textAlign: 'center',
                          padding: '4px'
                        }}>
                          {item.service.name}
                        </div>
                      ) : (
                        <img
                          src={item.service.image_url.startsWith('http') 
                            ? item.service.image_url 
                            : `http://localhost:3001/images/${item.service.image_url.replace(/^(img\/|images\/|deposits\/img\/)/, '')}`}
                          alt={item.service.name}
                          style={{
                            width: '80px',
                            height: '80px',
                            objectFit: 'cover',
                            borderRadius: '8px',
                            border: '1px solid #ddd'
                          }}
                          onError={(e) => {
                            e.currentTarget.style.display = 'none'
                          }}
                        />
                      )}
                    </div>
                  )}

                  {/* Детали вклада */}
                  <div style={{ flex: 1 }}>
                    <h3 style={{
                      color: '#5bd1d7',
                      fontSize: '20px',
                      margin: '0 0 8px 0',
                      fontWeight: '600'
                    }}>
                      Вклад «{item.service.name}»
                    </h3>
                    
                    {/* Ставка - только отображение (не редактируется) */}
                    <p style={{ margin: '5px 0', color: '#333', fontSize: '15px' }}>
                      Ставка: <strong>{item.service.price.toFixed(1)}%</strong>
                    </p>
                    
                    {/* Количество дней месяца - только отображение */}
                    <p style={{ margin: '5px 0', color: '#333', fontSize: '15px' }}>
                      Дней: <strong>{maxDays}</strong>
                    </p>

                    {/* Период - редактируемые поля для дат начала и окончания */}
                    <p style={{ margin: '5px 0', color: '#333', fontSize: '15px', display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>
                      Период:
                      <span style={{ fontSize: '14px', fontWeight: '500', color: '#47607d', marginRight: '6px' }}>
                        {item.service.name}
                      </span>
                      <input
                        type="number"
                        min="1"
                        max={maxDays}
                        value={editable.startDate}
                        onChange={(e) => handleStartDateChange(item.service_id, Number(e.target.value), maxDays)}
                        style={{
                          width: '60px',
                          padding: '4px 6px',
                          borderRadius: '6px',
                          border: '1px solid #c4d9f5',
                          fontSize: '14px',
                          textAlign: 'center',
                          backgroundColor: '#ffffff',
                          color: '#333333'
                        }}
                      />
                      <span style={{ margin: '0 4px' }}>—</span>
                      <input
                        type="number"
                        min="1"
                        max={maxDays}
                        value={editable.endDate}
                        onChange={(e) => handleEndDateChange(item.service_id, Number(e.target.value), maxDays)}
                        style={{
                          width: '60px',
                          padding: '4px 6px',
                          borderRadius: '6px',
                          border: '1px solid #c4d9f5',
                          fontSize: '14px',
                          textAlign: 'center',
                          backgroundColor: '#ffffff',
                          color: '#333333'
                        }}
                      />
                      <span style={{ fontSize: '14px', color: '#47607d', marginLeft: '4px' }}>числа</span>
                    </p>
                    
                    {/* Количество */}
                    <p style={{ margin: '5px 0', color: '#333', fontSize: '15px' }}>
                      Количество: <strong>{item.quantity}</strong>
                    </p>
                    
                    {/* Сумма вклада */}
                    <p style={{ margin: '5px 0', color: '#333', fontSize: '15px' }}>
                      Сумма вклада: <strong>{commonAmount.toLocaleString('ru-RU')} руб.</strong>
                    </p>
                    
                    {/* Прибыль - пересчитывается автоматически на основе периода */}
                    <p style={{ margin: '5px 0', color: '#333', fontSize: '15px' }}>
                      Прибыль: <strong>{profit.toFixed(2)} руб.</strong>
                    </p>
                  </div>

                  {/* Кнопка удаления */}
                  <div>
                    <button
                      type="button"
                      onClick={() => handleRemove(draft.id, item.service_id)}
                      style={{
                        backgroundColor: '#ff4d4f',
                        color: '#fff',
                        border: 'none',
                        padding: '8px 14px',
                        borderRadius: '6px',
                        fontSize: '14px',
                        cursor: 'pointer',
                        transition: '0.2s'
                      }}
                      onMouseOver={(e) => {
                        e.currentTarget.style.backgroundColor = '#e63946'
                      }}
                      onMouseOut={(e) => {
                        e.currentTarget.style.backgroundColor = '#ff4d4f'
                      }}
                    >
                      Удалить
                    </button>
                  </div>
                </div>
              )
            })}
          </div>

          {/* Кнопка отправки */}
          {draft.services.length > 0 && (
            <div style={{ textAlign: 'center', marginBottom: '30px' }}>
              <button
                type="button"
                onClick={handleSubmit}
                style={{
                  backgroundColor: '#009BFF',
                  color: '#fff',
                  border: 'none',
                  padding: '10px 18px',
                  borderRadius: '6px',
                  fontSize: '14px',
                  cursor: 'pointer',
                  transition: '0.2s'
                }}
                onMouseOver={(e) => {
                  e.currentTarget.style.backgroundColor = '#0077cc'
                }}
                onMouseOut={(e) => {
                  e.currentTarget.style.backgroundColor = '#009BFF'
                }}
              >
                Отправить расчет
              </button>
            </div>
          )}

          {/* Итоговая карточка */}
          {draft.services.length > 0 && (
            <section style={{
              marginTop: '30px',
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))',
              gap: '16px',
              background: '#ffffff',
              border: '1px solid #d6ecff',
              borderRadius: '12px',
              padding: '18px',
              boxShadow: '0 6px 20px rgba(0, 155, 255, 0.08)'
            }}>
              <div style={{
                background: '#f7fbff',
                border: '1px solid #e6f2ff',
                borderRadius: '10px',
                padding: '14px 16px'
              }}>
                <div style={{
                  fontSize: '14px',
                  color: '#47607d',
                  marginBottom: '6px',
                  fontWeight: '600'
                }}>
                  Общая сумма вкладов
                </div>
                <div style={{
                  fontSize: '22px',
                  fontWeight: '700',
                  color: '#0a2540'
                }}>
                  {totals.totalAmount.toLocaleString('ru-RU')} руб.
                </div>
              </div>
              <div style={{
                background: '#f7fbff',
                border: '1px solid #e6f2ff',
                borderRadius: '10px',
                padding: '14px 16px'
              }}>
                <div style={{
                  fontSize: '14px',
                  color: '#47607d',
                  marginBottom: '6px',
                  fontWeight: '600'
                }}>
                  Общая прибыль
                </div>
                <div style={{
                  fontSize: '22px',
                  fontWeight: '700',
                  color: '#0c7c3d'
                }}>
                  {totals.totalProfit.toFixed(2)} руб.
                </div>
              </div>
            </section>
          )}
        </>
      )}
    </div>
  )
}
