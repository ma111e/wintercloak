import {createApp} from 'vue'
import './assets/index.css'
import App from './App.vue'
import router from "../router.ts";
import { VCodeBlock } from '@wdns/vue-code-block';

const app = createApp(App);

app
    .component('VCodeBlock', VCodeBlock)
    .use(router)
    .mount('#app')
