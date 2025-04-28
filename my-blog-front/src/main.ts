import { createApp } from 'vue'
import App from './App.vue'

import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'virtual:uno.css'
import './styles/index.scss'

const app = createApp(App)

app.use(store)
app.use(router)
app.use(ElementPlus, { size: 'default', locale: zhCn })
app.mount('#app')
