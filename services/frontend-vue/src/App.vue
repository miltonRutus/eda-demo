<script setup>
import DashboardView from './views/DashboardView.vue';
</script>

<template>
  <v-app style="background: transparent;">
    <!-- Top Navigation Bar - MD3 Expressive Glassmorphic Header -->
    <v-app-bar
      elevation="0"
      density="comfortable"
      style="background: rgba(7, 10, 19, 0.75) !important; backdrop-filter: blur(20px); border-bottom: 1px solid rgba(255, 255, 255, 0.08);"
      class="px-2 px-md-4"
    >
      <template #prepend>
        <v-avatar
          size="40"
          class="mr-3"
          style="background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); box-shadow: 0 4px 14px rgba(99, 102, 241, 0.4);"
        >
          <v-icon size="22" color="white">mdi-lightning-bolt</v-icon>
        </v-avatar>
      </template>

      <v-app-bar-title>
        <div class="d-flex align-center">
          <span class="font-weight-black text-h6 text-white" style="letter-spacing: -0.02em;">EDA Platform</span>
          <v-chip
            size="x-small"
            color="primary"
            variant="tonal"
            class="ml-2 font-weight-bold d-none d-sm-inline-flex"
            style="border-radius: 6px !important;"
          >
            SOTA 2026
          </v-chip>
          <span class="text-caption text-grey-lighten-1 ml-2 font-weight-medium d-none d-md-inline">
            Event-Driven Architecture & Kong Gateway
          </span>
        </div>
      </v-app-bar-title>

      <template #append>
        <div class="d-flex align-center ga-2">
          <!-- Gateway Status Badge -->
          <v-chip
            size="small"
            color="primary"
            variant="tonal"
            class="px-3"
            style="border: 1px solid rgba(99, 102, 241, 0.3);"
          >
            <v-icon start size="14" color="primary">mdi-shield-check</v-icon>
            Kong DB-less Ingress
          </v-chip>

          <!-- OpenAPI Docs Links (Pill Buttons) -->
          <div class="d-none d-lg-flex align-center ga-1 ml-2 pl-2 border-s-grey-darken-3">
            <v-btn
              size="small"
              variant="tonal"
              color="surface-bright"
              prepend-icon="mdi-file-document-outline"
              href="http://localhost:8000/api/v1/fastapi/docs"
              target="_blank"
              class="md3-btn-expressive text-grey-lighten-1"
            >
              FastAPI Docs
            </v-btn>
            <v-btn
              size="small"
              variant="tonal"
              color="surface-bright"
              prepend-icon="mdi-file-document-outline"
              href="http://localhost:8000/api/v1/go/docs"
              target="_blank"
              class="md3-btn-expressive text-grey-lighten-1"
            >
              Go Docs
            </v-btn>
            <v-btn
              size="small"
              variant="tonal"
              color="surface-bright"
              prepend-icon="mdi-file-document-outline"
              href="http://localhost:8000/api/v1/legacy/docs"
              target="_blank"
              class="md3-btn-expressive text-grey-lighten-1"
            >
              Legacy Docs
            </v-btn>
          </div>
        </div>
      </template>
    </v-app-bar>

    <!-- Main Content Area with Ambient Glow Background -->
    <v-main style="background: transparent; min-height: 100vh;">
      <DashboardView />
    </v-main>
  </v-app>
</template>
