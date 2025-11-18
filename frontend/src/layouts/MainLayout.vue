<template>
  <a-layout class="layout">
    <a-layout-header class="header">
      <div class="logo">
        <book-outlined class="logo-icon" />
        <span class="logo-text">Post Management</span>
      </div>
      <a-menu
        v-model:selectedKeys="selectedKeys"
        theme="dark"
        mode="horizontal"
        :style="{ lineHeight: '64px' }"
      >
        <a-menu-item key="posts" @click="$router.push('/')">
          <unordered-list-outlined />
          <span>All Posts</span>
        </a-menu-item>
        <a-menu-item key="create" @click="$router.push('/posts/create')">
          <plus-circle-outlined />
          <span>Create Post</span>
        </a-menu-item>
      </a-menu>
    </a-layout-header>
    <a-layout-content class="content">
      <div class="content-wrapper">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </a-layout-content>
    <a-layout-footer class="footer">
      <div class="footer-content">
        Post Management System © {{ new Date().getFullYear() }}
      </div>
    </a-layout-footer>
  </a-layout>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import {
  BookOutlined,
  UnorderedListOutlined,
  PlusCircleOutlined,
} from '@ant-design/icons-vue'

const route = useRoute()
const selectedKeys = ref(['posts'])

// Update selected menu based on route
watch(
  () => route.path,
  (newPath) => {
    if (newPath === '/') {
      selectedKeys.value = ['posts']
    } else if (newPath.includes('/posts/create')) {
      selectedKeys.value = ['create']
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.layout {
  min-height: 100vh;
}

.header {
  display: flex;
  align-items: center;
  padding: 0 24px;
  background: #001529;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.logo {
  display: flex;
  align-items: center;
  margin-right: 40px;
  color: #fff;
  font-size: 20px;
  font-weight: 600;
  white-space: nowrap;
}

.logo-icon {
  font-size: 28px;
  margin-right: 12px;
  color: #1890ff;
}

.logo-text {
  letter-spacing: 0.5px;
}

.content {
  background: #f0f2f5;
  padding: 24px;
}

.content-wrapper {
  max-width: 1400px;
  margin: 0 auto;
}

.footer {
  text-align: center;
  background: #fff;
  border-top: 1px solid #f0f0f0;
}

.footer-content {
  color: rgba(0, 0, 0, 0.45);
}

/* Page transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
