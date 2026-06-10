<template>
  <div id="app">
    <AppNav v-if="!route.meta?.hideNav" />
    <router-view v-slot="{ Component, route: currentRoute }">
      <transition name="page" mode="out-in">
        <component :is="Component" :key="currentRoute.path" />
      </transition>
    </router-view>
    <AppFooter v-if="!route.meta?.hideNav" />
    <SearchModal ref="searchRef" />
  </div>
</template>

<script setup>
import { ref, provide } from 'vue'
import { useRoute } from 'vue-router'
import AppNav from '@/components/AppNav.vue'
import AppFooter from '@/components/AppFooter.vue'
import SearchModal from '@/components/SearchModal.vue'

const route = useRoute()
const searchRef = ref(null)

provide('openSearch', () => {
  searchRef.value?.openModal()
})
</script>

<style>
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
</style>
