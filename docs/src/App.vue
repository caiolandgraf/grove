<template>
  <div id="app">
    <AppNav />
    <router-view v-slot="{ Component, route }">
      <transition name="page" mode="out-in">
        <component :is="Component" :key="route.path" />
      </transition>
    </router-view>
    <AppFooter />
    <SearchModal ref="searchRef" />
  </div>
</template>

<script setup>
import { ref, provide } from 'vue'
import AppNav from '@/components/AppNav.vue'
import AppFooter from '@/components/AppFooter.vue'
import SearchModal from '@/components/SearchModal.vue'

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
