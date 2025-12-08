# 🔧 ИСПРАВЛЕНИЕ ПУСТОЙ СТРАНИЦЫ НА GITHUB PAGES

## 🔍 ДИАГНОСТИКА

Если на `akabracom.github.io/zi/` белый экран, проверьте следующее:

---

## 1️⃣ ПРОВЕРЬТЕ WORKFLOW

1. Откройте: https://github.com/akabracom/zi/actions
2. Найдите последний запуск workflow для ветки `tauri-deposits`
3. Проверьте:
   - ✅ **Зеленая галочка** = workflow выполнился успешно
   - ❌ **Красный крестик** = есть ошибки (откройте и посмотрите логи)
   - ⏳ **Желтый кружок** = еще выполняется (подождите)

**Если workflow не запустился:**
- Убедитесь, что файл `.github/workflows/deploy.yml` есть в ветке `tauri-deposits`
- Убедитесь, что вы запушили этот файл

---

## 2️⃣ ПРОВЕРЬТЕ НАСТРОЙКИ PAGES

1. Откройте: https://github.com/akabracom/zi/settings/pages
2. Убедитесь, что:
   - **Source:** `Deploy from a branch`
   - **Branch:** `tauri-deposits` (или `spa-deposits` для 5 лабы)
   - **Folder:** `/ (root)`
3. Нажмите **Save** (даже если уже сохранено)

---

## 3️⃣ ПРОВЕРЬТЕ КОНФИГУРАЦИЮ

### Для 6 лабы (`tauri-deposits`):

**`vite.config.ts` должен содержать:**
```typescript
export default defineConfig({
  base: '/zi/', // ОБЯЗАТЕЛЬНО!
  plugins: [react()],
  // ...
})
```

**`src/App.tsx` должен содержать:**
```typescript
<BrowserRouter basename={import.meta.env.BASE_URL || '/zi/'}>
```

### Для 5 лабы (`spa-deposits`):

**`vite.config.ts` должен быть БЕЗ `base`:**
```typescript
export default defineConfig({
  // БЕЗ base!
  plugins: [react()],
  // ...
})
```

**`src/App.tsx` должен быть БЕЗ `basename`:**
```typescript
<BrowserRouter>
```

---

## 4️⃣ ИСПРАВЬТЕ WORKFLOW (если нужно)

Если workflow падает с ошибкой, обновите `.github/workflows/deploy.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches:
      - tauri-deposits  # или spa-deposits для 5 лабы

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

jobs:
  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      
      - name: Install dependencies
        run: npm install  # Изменили с npm ci на npm install
      
      - name: Build
        run: npm run build
        env:
          VITE_USE_MOCK: true
      
      - name: Setup Pages
        uses: actions/configure-pages@v4
      
      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: './dist'
      
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

**Изменение:** `npm ci` → `npm install` (если нет `package-lock.json`)

---

## 5️⃣ ПЕРЕЗАПУСТИТЕ DEPLOYMENT

1. Откройте: https://github.com/akabracom/zi/actions
2. Найдите последний workflow
3. Нажмите **Re-run all jobs** (если есть ошибки)
4. Или сделайте пустой коммит и запушьте:
   ```bash
   git commit --allow-empty -m "Trigger deployment"
   git push origin tauri-deposits
   ```

---

## 6️⃣ ПРОВЕРЬТЕ КОНСОЛЬ БРАУЗЕРА

1. Откройте `akabracom.github.io/zi/`
2. Нажмите **F12** (открыть DevTools)
3. Перейдите на вкладку **Console**
4. Посмотрите на ошибки (красные сообщения)

**Частые ошибки:**
- `404` на файлы `.js` или `.css` → проблема с `base` в `vite.config.ts`
- `CORS error` → нормально для GitHub Pages (используйте mock режим)
- `Cannot find module` → проблема со сборкой

---

## 7️⃣ ПРОВЕРЬТЕ СЕТЬ

1. В DevTools перейдите на вкладку **Network**
2. Обновите страницу (F5)
3. Проверьте:
   - Загружается ли `index.html`? (должен быть статус 200)
   - Загружаются ли файлы `.js` и `.css`? (должны быть статус 200)

**Если файлы не загружаются:**
- Проверьте пути в `index.html` (должны начинаться с `/zi/`)
- Проверьте `base` в `vite.config.ts`

---

## 8️⃣ ОЧИСТИТЕ КЭШ

1. **В браузере:**
   - Ctrl+Shift+Delete → очистить кэш
   - Или Ctrl+Shift+R (жесткая перезагрузка)

2. **На GitHub:**
   - Settings → Pages → Clear cache (если есть)

---

## ✅ ЧЕКЛИСТ

- [ ] Workflow выполнился успешно (зеленая галочка)
- [ ] Настройки Pages правильные (ветка и папка)
- [ ] `base: '/zi/'` в `vite.config.ts` (для 6 лабы)
- [ ] `basename={import.meta.env.BASE_URL || '/zi/'}` в `App.tsx` (для 6 лабы)
- [ ] Нет ошибок в консоли браузера
- [ ] Файлы загружаются в Network tab

---

## 🚀 БЫСТРОЕ ИСПРАВЛЕНИЕ

Если ничего не помогло, выполните:

```bash
# 1. Перейти в папку проекта
cd C:\spa-frontend\lab6\spa-deposits-tauri

# 2. Проверить конфигурацию
# Убедитесь, что vite.config.ts содержит base: '/zi/'

# 3. Сделать пустой коммит для перезапуска
git commit --allow-empty -m "Fix GitHub Pages deployment"
git push origin tauri-deposits

# 4. Подождать 2-3 минуты
# 5. Проверить: https://github.com/akabracom/zi/actions
# 6. Открыть: https://akabracom.github.io/zi/
```

---

**После исправления сайт должен загрузиться!** 🎉

