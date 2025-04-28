import { $get, $post, $delete } from '@/utils/request'

export const ajaxArticle = {
  query: (data: any) => $post('/common/article/query', data),
  add: (data: any) => $post('/article/add', data),
  getById: (id: number) => $get(`/common/article/getById/${ id }`),
  getMyArticle: (data: any) => $post('/article/getMyArticle', data),
  deleteArticle: (id: number) => $delete(`/article/delete/${ id }`),

  queryReply: (data: any) => $post('/common/reply/query', data),
  addReply: (data: any) => $post('/reply/add', data),

  addLike: (data: { articleId: number }) => $post('/like/add', data),
  cancelLike: (data: { articleId: number }) => $post('/like/cancel', data)
}

export const ajaxReply = {
  query: (data: any) => $post('/reply/query', data),
  add: (data: any) => $post('/reply/add', data),
}
