<template>
  <div class="main_Container details h-[calc(100vh-80px)] overflow-auto pt-10px bg-[var(--bg-color-2)]">
    <blog-detail />
  </div>
</template>

<script setup lang="ts">
import { ref, provide, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import blogDetail from '@/components/blogDetail/index.vue'

const route = useRoute()
const router = useRouter()
const id = ref(route.params.id)

provide('id', id)

onMounted(() => {
  const numId = Number(id.value)
  if (numId === 0 || numId + '' === 'NaN') {
    router.push('/')
  }
})
</script>

<style scoped lang="scss">
.details {
  background-clip: content-box;
}
</style>
