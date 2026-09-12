<script setup>
import DashboardView from './views/DashboardView.vue';
</script>

<template>
  <v-app>
    <!-- Top Navigation Bar -->
    <v-app-bar color="#090d16" elevation="1" density="comfortable" class="border-b-grey-darken-3">
      <template #prepend>
        <v-avatar color="primary" size="36" class="mr-2">
          <v-icon size="20" color="white">mdi-lightning-bolt</v-icon>
        </v-avatar>
      </template>

      <v-app-bar-title class="font-weight-bold text-subtitle-1">
        EDA Platform <span class="text-caption text-grey ml-2 font-weight-regular d-none d-sm-inline">Event-Driven Architecture & Kong Gateway</span>
      </v-app-bar-title>

      <template #append>
        <div class="d-flex align-center ga-2 mr-2">
          <!-- Gateway Status Badge -->
          <v-chip size="small" color="primary" variant="outlined">
            <v-icon start size="14">mdi-shield-check</v-icon>
            Kong DB-less
          </v-chip>

          <!-- OpenAPI Docs Links -->
          <v-btn
            size="small"
            variant="text"
            prepend-icon="mdi-file-document-outline"
            href="http://localhost:8000/api/v1/fastapi/docs"
            target="_blank"
            class="d-none d-md-inline-flex"
          >
            FastAPI Docs
          </v-btn>
          <v-btn
            size="small"
            variant="text"
            prepend-icon="mdi-file-document-outline"
            href="http://localhost:8000/api/v1/go/docs"
            target="_blank"
            class="d-none d-md-inline-flex"
          >
            Go Docs
          </v-btn>
          <v-btn
            size="small"
            variant="text"
            prepend-icon="mdi-file-document-outline"
            href="http://localhost:8000/api/v1/legacy/docs"
            target="_blank"
            class="d-none d-md-inline-flex"
          >
            Legacy Docs
          </v-btn>
        </div>
      </template>
    </v-app-bar>

    <!-- Main Content Area -->
    <v-main style="background-color: #0a0e17; min-height: 100vh;">
      <DashboardView />
    </v-main>
  </v-app>
</template>
