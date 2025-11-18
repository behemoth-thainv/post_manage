import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'post-list',
        component: () => import('@/views/PostList.vue'),
        meta: {
          title: 'All Posts',
        },
      },
      {
        path: '/posts/create',
        name: 'post-create',
        component: () => import('@/views/PostForm.vue'),
        meta: {
          title: 'Create Post',
        },
      },
      {
        path: '/posts/:id',
        name: 'post-detail',
        component: () => import('@/views/PostDetail.vue'),
        meta: {
          title: 'Post Details',
        },
      },
      {
        path: '/posts/:id/edit',
        name: 'post-edit',
        component: () => import('@/views/PostForm.vue'),
        meta: {
          title: 'Edit Post',
        },
      },
    ],
  },
  {
    // Catch all 404
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
})

// Navigation guard to update page title
router.beforeEach((to, from, next) => {
  document.title = to.meta.title
    ? `${to.meta.title} | Post Management`
    : 'Post Management System'
  next()
})

export default router
