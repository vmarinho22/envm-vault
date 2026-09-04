import { createPinia } from 'pinia';
import { createApp } from 'vue';
import App from '@/app/App.vue';
import router from '@/app/router';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { httpClient } from '@/shared/http';
import '@/app/styles/main.css';

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
httpClient.configure(() => useAuthStore(pinia).token);
app.use(router);

app.mount('#app');
