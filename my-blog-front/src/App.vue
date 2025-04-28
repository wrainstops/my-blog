<template>
  <div>
    <wave />
    <canvas
      ref="live2dRef"
      id="live2d"
      class="fixed h-200px w-auto right-0 bottom-0 z-999 pointer-events-none"
    />
    <el-tooltip
      :content="showLive2d ? '隐藏人物' : '显示人物'"
      effect="customized"
      placement="left"
    >
      <div
        id="showBtn"
        class="showBtn fixed w-40px h-40px right-150px bottom-200px border-rd-50% flex items-center justify-center z-5 pointer"
        @click="handleChangeShowLive2d"
      >
        {{ showLive2d ? '隐' : '显' }}
      </div>
    </el-tooltip>
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import * as PIXI from 'pixi.js'
import { InternalModel, Live2DModel } from 'pixi-live2d-display/cubism4'
import wave from '@/components/wave.vue'

window.PIXI = PIXI
let pixi: PIXI.Application
let model: Live2DModel<InternalModel>
const live2dRef = ref()

async function initPIXI() {
  pixi = new PIXI.Application({
    view: live2dRef.value,
    autoStart: true,
    backgroundAlpha: 0 // 透明度
  })
  model = await Live2DModel.from('/live2d/Ava/Ava.model3.json')
  pixi.stage.addChild(model)
  model.scale.set(0.5) // 整体缩放0.8
}

const showLive2d = ref(true)
function handleChangeShowLive2d() {
  showLive2d.value = !showLive2d.value
  if (showLive2d.value) {
    model.visible = true
  } else {
    model.visible = false
  }
}

onMounted(async () => {
  await initPIXI()
})
onUnmounted(() => {
  model?.destroy()
  pixi?.destroy()
})
</script>

<style scoped lang="scss">
.showBtn {
  background-color: var(--bg-color-1);
  box-shadow: var(--el-box-shadow-lighter);
}
.showBtn:hover {
  background-color: var(--bg-color-2);
}
</style>
