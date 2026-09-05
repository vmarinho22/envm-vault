import { createPinia } from 'pinia';
import { createApp } from 'vue';
import App from '@/app/App.vue';
import router from '@/app/router';
import { applyDocumentLang, i18n } from '@/app/translations';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { assertEnvsValid } from '@/shared/env/envs';
import { httpClient } from '@/shared/http';
import '@/app/styles/main.css';

assertEnvsValid();

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
httpClient.configure(() => useAuthStore(pinia).token);
app.use(i18n);
applyDocumentLang(i18n.global.locale.value);
app.use(router);

app.mount('#app');
