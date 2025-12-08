# 🚀 НАСТРОЙКА GITHUB PAGES

## 📋 ЧТО НУЖНО СДЕЛАТЬ

После того, как вы запушили ветку `tauri-deposits`, нужно настроить GitHub Pages.

---

## 🔧 ДЛЯ 6 ЛАБЫ (ветка `tauri-deposits`)

### 1. Настройка GitHub Pages

1. Откройте: https://github.com/akabracom/zi
2. Переключитесь на ветку `tauri-deposits`
3. Нажмите **Settings** → **Pages**
4. В разделе **Source**:
   - Выберите **Deploy from a branch**
   - Branch: `tauri-deposits`
   - Folder: `/ (root)`
5. Нажмите **Save**

### 2. Создание GitHub Actions Workflow

1. В ветке `tauri-deposits` создайте папку `.github/workflows` (если её нет)
2. Создайте файл `.github/workflows/deploy.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches:
      - tauri-deposits

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
        run: npm ci
      
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

3. Закоммитьте и запушьте:
```bash
git add .github/workflows/deploy.yml
git commit -m "Add GitHub Pages workflow"
git push origin tauri-deposits
```

### 3. URL после деплоя

После успешного деплоя ваш сайт будет доступен по адресу:

**https://akabracom.github.io/zi/**

---

## ⚙️ ВАЖНЫЕ НАСТРОЙКИ

### Проверьте `vite.config.ts`

**Для 6 лабы** (`tauri-deposits`):
```typescript
export default defineConfig({
  base: '/zi/', // ДЛЯ GitHub Pages
  plugins: [react()],
  server: { ... }
})
```

### Проверьте `App.tsx`

**Для 6 лабы** (`tauri-deposits`):
```typescript
<BrowserRouter basename={import.meta.env.BASE_URL || '/zi/'}>
  {/* С basename для GitHub Pages */}
```

---

## 📱 КАК ОТКРЫТЬ НА ТЕЛЕФОНЕ

1. Откройте браузер на телефоне
2. Перейдите по адресу: **https://akabracom.github.io/zi/**
3. Дождитесь загрузки сайта
4. Сохраните как PWA:
   - **Android (Chrome):** Меню (⋮) → "Добавить на главный экран"
   - **iOS (Safari):** Поделиться (□↑) → "На экран «Домой»"

---

## ✅ ПРОВЕРКА

После деплоя:

1. Откройте: https://github.com/akabracom/zi/actions
2. Убедитесь, что workflow выполнился успешно (зеленая галочка)
3. Откройте: https://akabracom.github.io/zi/
4. Проверьте, что сайт загружается

---

## 🐛 ЕСЛИ НЕ РАБОТАЕТ

1. **Проверьте workflow:**
   - Откройте: https://github.com/akabracom/zi/actions
   - Посмотрите логи ошибок

2. **Проверьте настройки Pages:**
   - Settings → Pages
   - Убедитесь, что выбрана правильная ветка

3. **Проверьте `base` в `vite.config.ts`:**
   - Для 6 лабы: `base: '/zi/'`

4. **Очистите кэш браузера:**
   - Ctrl+Shift+R (Windows) или Cmd+Shift+R (Mac)

---

**Готово! После настройки ветка будет доступна через GitHub Pages.** 🎉

