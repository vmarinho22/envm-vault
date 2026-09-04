import { createPinia } from 'pinia';
import { createApp } from 'vue';
import App from '@/app/App.vue';
import router from '@/app/router';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { assertEnvsValid } from '@/shared/env/envs';
import { httpClient } from '@/shared/http';
import '@/app/styles/main.css';

assertEnvsValid();

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
httpClient.configure(() => useAuthStore(pinia).token);
app.use(router);

app.mount('#app');
