---
name: add-ui-screen
description: Checklist for adding a new Vue screen or view in the envm-vault admin UI. Use when creating pages, routes, or UI features in the frontend.
---

# Add UI Screen

Use this checklist when adding a new screen to the Vue admin UI.

## Before coding

1. Read [docs/architecture.md](../../docs/architecture.md) for domain context.
2. Identify which API endpoints the screen needs.
3. Choose or create the module under `frontend/src/modules/{domain}/`.
4. Check if a service, composable, or store already exists in that module.

## Implementation checklist

```
- [ ] Add or extend a typed service: `interface XService` + `export const xService: XService`
- [ ] Create view in modules/{domain}/views/
- [ ] Register route in app/router/index.ts
- [ ] Create or extend Pinia store in the module if shared state is needed
- [ ] Create composable in modules/{domain}/composables/ for UI flow
- [ ] Handle loading, error, and empty states
- [ ] Add Vitest coverage: unit under `tests/unit/` (mirrors module), components under `tests/component/` (mirrors src path)
- [ ] Every component test calls `expectComponentSnapshot(wrapper)` (commit the generated `__snapshots__`)
- [ ] Add nav item in app/layouts/DefaultLayout.vue if needed
```

## View pattern

```vue
<script setup lang="ts">
import { useProjects } from '@/modules/project/composables/useProjects'

const { projects, loading, error, fetchProjects } = useProjects()

onMounted(() => fetchProjects())
</script>

<template>
  <div>
    <p v-if="loading">Carregando...</p>
    <p v-else-if="error" role="alert">{{ error }}</p>
    <ProjectList v-else :projects="projects" />
  </div>
</template>
```

## Component test pattern

```ts
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount'

it('renders …', async () => {
  const { wrapper } = await mountComponent(MyView)
  // assertions…
  expectComponentSnapshot(wrapper)
})
```

Commit the generated `__snapshots__/` files. Update with `npm run test:component -- -u` only when the UI change is intentional.

## Standards

- Use `<script setup lang="ts">` — no Options API
- Extract logic into module composables, not inline in views
- Call the API only through typed service instances (`projectService.list()`, etc.)
- Never store decrypted secret values in Pinia beyond the current view lifecycle
- Include accessible labels on form inputs
- Every component test calls `expectComponentSnapshot(wrapper)`
