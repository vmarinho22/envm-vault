---
name: add-ui-screen
description: Checklist for adding a new Vue screen or view in the envm-vault admin UI. Use when creating pages, routes, or UI features in the frontend.
---

# Add UI Screen

Use this checklist when adding a new screen to the Vue admin UI (Phase 3+).

## Before coding

1. Read [docs/architecture.md](../../docs/architecture.md) for domain context.
2. Identify which API endpoints the screen needs.
3. Check if a composable or store already exists for the data.

## Implementation checklist

```
- [ ] Create view component in src/views/
- [ ] Add route in router config
- [ ] Create or extend Pinia store if shared state is needed
- [ ] Create composable in src/composables/ for data fetching
- [ ] Add API client function in src/api/
- [ ] Handle loading, error, and empty states
- [ ] Add navigation link if needed
```

## View pattern

```vue
<script setup lang="ts">
import { useProjects } from '@/composables/useProjects'

const { projects, loading, error, fetchProjects } = useProjects()

onMounted(() => fetchProjects())
</script>

<template>
  <div>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error">{{ error }}</p>
    <ProjectList v-else :projects="projects" />
  </div>
</template>
```

## Standards

- Use `<script setup lang="ts">` — no Options API
- Extract logic into composables, not inline in views
- Never store decrypted secret values in Pinia beyond the current view lifecycle
- Include accessible labels on form inputs
