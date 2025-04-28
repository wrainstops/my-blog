<template>
  <div class="page404">
    <el-image class="img" :src="Img404" fit="contain"></el-image>
  </div>
</template>

<script setup>
import Img404 from '@/assets/image/404.png'
</script>

<style scoped>
.page404 {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #DBDBDC;

  .img {
    width: 50%;
    height: 50%;
  }
}
</style>
