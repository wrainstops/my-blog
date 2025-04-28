<template>
  <div class="tinymce" :style="{'--color': backColor }">
    <editor
      :init="init"
      :disabled="props.disabled"
    ></editor>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { isDark } from '@/composables/dark'

import tinymce from 'tinymce/tinymce' // tinymce默认hidden，不引入不显示
import Editor from '@tinymce/tinymce-vue'

// 引入富文本编辑器主题的js和css
import 'tinymce/skins/content/default/content.css'
import 'tinymce/themes/silver/theme.min.js'
import 'tinymce/icons/default/icons' // 解决了icons.js 报错Unexpected token '<'

import 'tinymce/plugins/image' // 插入上传图片插件
import 'tinymce/plugins/media' // 插入视频插件
import 'tinymce/plugins/table' // 插入表格插件
import 'tinymce/plugins/lists' // 列表插件
import 'tinymce/plugins/wordcount' // 字数统计插件
import 'tinymce/plugins/link'
import 'tinymce/plugins/code'
import 'tinymce/plugins/preview'
import 'tinymce/plugins/fullscreen'
import 'tinymce/models/dom/model'

const props = reactive({
  // 默认的富文本内容
  value: '',
  // 基本路径，默认为空根目录，如果你的项目发布后的地址为目录形式，
  // 即abc.com/tinymce，baseUrl需要配置成tinymce，不然发布后资源会找不到
  baseUrl: window.location.origin ? window.location.origin : '',
  disabled: false,
  plugins: 'link lists image code table wordcount media preview fullscreen',
  toolbar: 'bold italic underline strikethrough | fontsizeselect | formatselect | forecolor backcolor | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent blockquote | undo redo | link unlink code lists table image media | removeformat | fullscreen preview'
})
const init = reactive({
  language_url: `${ props.baseUrl }/tinymce/langs/zh-Hans.js`,
  language: 'zh-Hans',
  // skin_url: `${ props.baseUrl }/tinymce/skins/ui/oxide`,
  skin_url: isDark.value ? `${ props.baseUrl }/tinymce/skins/ui/oxide-dark` : `${ props.baseUrl }/tinymce/skins/ui/oxide`, // 暗色系
  convert_urls: false,
  height: '70vh',
  // content_css（为编辑区指定css文件）,加上就不显示字数统计了
  content_css: '',
  // 指定需加载的插件
  plugins: props.plugins,
  // 自定义工具栏
  toolbar: props.toolbar,
  // 显示顶部升级按钮
  promotion: false,
  // 底部的状态栏
  statusbar: true,
  // 最上方的菜单
  menubar: 'file edit insert view format table tools help',
  // 显示右下角水印 “Powered by TinyMCE”
  branding: false,
  // 此处为图片上传处理函数，这个直接用了base64的图片形式上传图片，
  // 如需ajax上传可参考https://www.tiny.cloud/docs/configure/file-image-upload/#images_upload_handler
  images_upload_handler: (blobInfo: { base64: () => string }, success: (arg0: string) => void, failure: any) => {
    const img = 'data:image/jpeg;base64,' + blobInfo.base64()
    success(img)
    console.log(failure)
  }
})

tinymce.init({})

const backColor = ref(isDark.value ? '#424171' : '#ddd')
watch(
  () => isDark.value,
  (val: boolean) => {
    if (val) {
      backColor.value = '#424171'
      init.skin_url = `${ props.baseUrl }/tinymce/skins/ui/oxide-dark`
    } else {
      backColor.value = '#ddd'
      init.skin_url = `${ props.baseUrl }/tinymce/skins/ui/oxide`
    }
  }
)
defineExpose({
  tinymce
})
</script>

<style scoped lang="scss">
:deep(.tox .tox-edit-area__iframe) {
  background-color: var(--color);
}
:deep(.tox-tinymce) {
  border: 0;
}
</style>
