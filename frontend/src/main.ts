import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { initTheme } from './lib/theme'
import 'vue-sonner/style.css'
import './style.css'

initTheme()

createApp(App).use(router).mount('#app')
