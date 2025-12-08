import { useEffect, useState } from 'react'
import { completeRequest, deleteRequestLogical, listRequests } from '../api/http'
import type { Request } from '../api/types'
import { useRole } from '../context/AuthContext'

export function RequestsPage() {
  const [requests, setRequests] = useState<Request[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const role = useRole()

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const data = await listRequests()
      setRequests(data)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Не удалось загрузить расчеты',
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void load()
  }, [])

  async function handleComplete(id: number) {
    try {
      await completeRequest(id)
      await load()
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Не удалось завершить заявку',
      )
    }
  }

  async function handleDelete(id: number) {
    try {
      await deleteRequestLogical(id)
      await load()
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Не удалось удалить заявку',
      )
    }
  }

  return (
    <div className="page">
      <h2>Расчеты вкладов</h2>
      {error && <div className="form-error">{error}</div>}
      {loading && <p>Загрузка...</p>}

      <table className="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Статус</th>
            <th>Дата создания</th>
            <th>Количество месяцев</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {requests.map((r) => (
            <tr key={r.id}>
              <td>{r.id}</td>
              <td>{r.status}</td>
              <td>{new Date(r.created_at).toLocaleString()}</td>
              <td>{r.services.length}</td>
              <td>
                {role === 'moderator' && (
                  <>
                    <button
                      type="button"
                      className="link-button"
                      onClick={() => handleComplete(r.id)}
                    >
                      Завершить
                    </button>
                    <button
                      type="button"
                      className="link-button"
                      onClick={() => handleDelete(r.id)}
                    >
                      Удалить
                    </button>
                  </>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {!loading && requests.length === 0 && (
        <p className="muted">Расчетов пока нет.</p>
      )}
    </div>
  )
}


