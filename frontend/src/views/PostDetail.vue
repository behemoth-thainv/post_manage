<template>
  <div class="post-detail">
    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <a-spin size="large" tip="Loading post details..." />
    </div>

    <!-- Post Content -->
    <template v-else-if="post">
      <!-- Header -->
      <a-page-header
        class="page-header"
        :title="post.title"
        @back="$router.back()"
      >
        <template #extra>
          <a-space>
            <a-button @click="$router.push(`/posts/${post.id}/edit`)">
              <template #icon>
                <edit-outlined />
              </template>
              Edit
            </a-button>
            <a-popconfirm
              title="Are you sure you want to delete this post?"
              ok-text="Yes"
              cancel-text="No"
              @confirm="handleDelete"
            >
              <a-button danger>
                <template #icon>
                  <delete-outlined />
                </template>
                Delete
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>

        <template #footer>
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="Created">
              <calendar-outlined />
              {{ formatDate(post.created_at) }}
            </a-descriptions-item>
            <a-descriptions-item label="Last Updated">
              <clock-circle-outlined />
              {{ formatDate(post.updated_at) }}
            </a-descriptions-item>
          </a-descriptions>
        </template>
      </a-page-header>

      <!-- Content Card -->
      <a-card class="content-card" :bordered="false">
        <div class="post-content" v-html="formatContent(post.content)"></div>
      </a-card>

      <!-- Action Card -->
      <a-card class="action-card" :bordered="false">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-button
              type="primary"
              size="large"
              block
              @click="$router.push(`/posts/${post.id}/edit`)"
            >
              <template #icon>
                <edit-outlined />
              </template>
              Edit Post
            </a-button>
          </a-col>
          <a-col :span="12">
            <a-button size="large" block @click="$router.push('/')">
              <template #icon>
                <arrow-left-outlined />
              </template>
              Back to List
            </a-button>
          </a-col>
        </a-row>
      </a-card>
    </template>

    <!-- Error State -->
    <a-empty v-else description="Post not found">
      <a-button type="primary" @click="$router.push('/')">
        Go to Posts List
      </a-button>
    </a-empty>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
import {
  EditOutlined,
  DeleteOutlined,
  CalendarOutlined,
  ClockCircleOutlined,
  ArrowLeftOutlined,
} from '@ant-design/icons-vue'
import { postService } from '@/services/api'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const post = ref(null)

// Load post details
const loadPost = async () => {
  loading.value = true
  try {
    const response = await postService.getPost(route.params.id)
    post.value = response.data
  } catch (error) {
    message.error('Failed to load post')
    console.error('Error loading post:', error)
  } finally {
    loading.value = false
  }
}

// Handle delete
const handleDelete = async () => {
  try {
    await postService.deletePost(route.params.id)
    message.success('Post deleted successfully')
    router.push('/')
  } catch (error) {
    message.error('Failed to delete post')
    console.error('Error deleting post:', error)
  }
}

// Format date
const formatDate = (date) => {
  return dayjs(date).format('MMMM D, YYYY [at] h:mm A')
}

// Format content with line breaks
const formatContent = (content) => {
  if (!content) return ''
  return content.replace(/\n/g, '<br>')
}

// Initialize
onMounted(() => {
  loadPost()
})
</script>

<style scoped>
.post-detail {
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.page-header {
  background: #fff;
  margin-bottom: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.page-header :deep(.ant-page-header-heading-title) {
  font-size: 32px;
  font-weight: 700;
  color: #262626;
}

.content-card {
  margin-bottom: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  min-height: 300px;
}

.post-content {
  font-size: 16px;
  line-height: 1.8;
  color: rgba(0, 0, 0, 0.85);
  white-space: pre-wrap;
  word-wrap: break-word;
}

.action-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}
</style>
