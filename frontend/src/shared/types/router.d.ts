export {};

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean;
    requiresApiToken?: boolean;
    requiresAdmin?: boolean;
  }
}
