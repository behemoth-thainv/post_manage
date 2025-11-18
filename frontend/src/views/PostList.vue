<template>
  <div class="post-list">
    <!-- Page Header -->
    <a-page-header
      class="page-header"
      title="All Posts"
      sub-title="Manage your blog posts"
    >
      <template #extra>
        <a-button type="primary" size="large" @click="$router.push('/posts/create')">
          <template #icon>
            <plus-outlined />
          </template>
          Create New Post
        </a-button>
      </template>
    </a-page-header>

    <!-- Search and Filter -->
    <a-card class="filter-card" :bordered="false">
      <a-row :gutter="16">
        <a-col :span="18">
          <a-input-search
            v-model:value="searchText"
            placeholder="Search posts by title or content..."
            size="large"
            @search="handleSearch"
            allow-clear
          >
            <template #prefix>
              <search-outlined />
            </template>
          </a-input-search>
        </a-col>
        <a-col :span="6">
          <a-button size="large" block @click="loadPosts">
            <template #icon>
              <reload-outlined />
            </template>
            Refresh
          </a-button>
        </a-col>
      </a-row>
    </a-card>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <a-spin size="large" tip="Loading posts...">
        <div style="padding: 50px" />
      </a-spin>
    </div>

    <!-- Empty State -->
    <a-empty
      v-else-if="!loading && filteredPosts.length === 0"
      description="No posts found"
      class="empty-state"
    >
      <a-button type="primary" @click="$router.push('/posts/create')">
        Create First Post
      </a-button>
    </a-empty>

    <!-- Posts Grid -->
    <a-row v-else :gutter="[24, 24]" class="posts-grid">
      <a-col
        v-for="post in filteredPosts"
        :key="post.id"
        :xs="24"
        :sm="24"
        :md="12"
        :lg="8"
        :xl="8"
      >
        <a-card
          hoverable
          class="post-card"
          :body-style="{ padding: '24px' }"
        >
          <template #actions>
            <a-tooltip title="View Details">
              <eye-outlined @click="viewPost(post.id)" />
            </a-tooltip>
            <a-tooltip title="Edit Post">
              <edit-outlined @click="editPost(post.id)" />
            </a-tooltip>
            <a-tooltip title="Delete Post">
              <a-popconfirm
                title="Are you sure you want to delete this post?"
                ok-text="Yes"
                cancel-text="No"
                @confirm="deletePost(post.id)"
              >
                <delete-outlined style="color: #ff4d4f" />
              </a-popconfirm>
            </a-tooltip>
          </template>

          <a-card-meta>
            <template #title>
              <div class="post-title" @click="viewPost(post.id)">
                {{ post.title }}
              </div>
            </template>
            <template #description>
              <div class="post-content">
                {{ truncateContent(post.content) }}
              </div>
              <div class="post-meta">
                <a-space :size="16">
                  <span>
                    <calendar-outlined />
                    {{ formatDate(post.created_at) }}
                  </span>
                  <span v-if="post.updated_at !== post.created_at">
                    <clock-circle-outlined />
                    Updated {{ formatDate(post.updated_at) }}
                  </span>
                </a-space>
              </div>
            </template>
          </a-card-meta>
        </a-card>
      </a-col>
    </a-row>

    <!-- Pagination -->
    <a-card v-if="!loading && !searchText && pagination.count > 0" class="pagination-card" :bordered="false">
      <a-pagination
        v-model:current="currentPage"
        v-model:page-size="pageSize"
        :total="pagination.count"
        :show-total="(total, range) => `${range[0]}-${range[1]} of ${total} posts`"
        :show-size-changer="true"
        :page-size-options="['12', '24', '36', '48']"
        @change="handlePageChange"
        @show-size-change="handlePageSizeChange"
        :show-quick-jumper="true"
        size="default"
        class="custom-pagination"
      />
    </a-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import {
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
  EyeOutlined,
  EditOutlined,
  DeleteOutlined,
  CalendarOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons-vue'
import { postService } from '@/services/api'

dayjs.extend(relativeTime)

const router = useRouter()
const loading = ref(false)
const posts = ref([])
const searchText = ref('')
const currentPage = ref(1)
const pageSize = ref(12)
const pagination = ref({
  page: 1,
  limit: 12,
  count: 0,
  pages: 0,
  prev_page: null,
  next_page: null
})

// Computed property for filtered posts (only for search)
const filteredPosts = computed(() => {
  if (!searchText.value) {
    return posts.value
  }
  const search = searchText.value.toLowerCase()
  return posts.value.filter(
    (post) =>
      post.title.toLowerCase().includes(search) ||
      post.content.toLowerCase().includes(search)
  )
})

// Load posts with pagination
const loadPosts = async (page = currentPage.value, limit = pageSize.value) => {
  loading.value = true
  try {
    const response = await postService.getPosts(page, limit)
    posts.value = response.data.posts
    pagination.value = response.data.pagination
    
    // Update current page and page size
    currentPage.value = pagination.value.page
    pageSize.value = pagination.value.limit
    
    message.success(`Loaded ${posts.value.length} posts`)
  } catch (error) {
    message.error('Failed to load posts')
    console.error('Error loading posts:', error)
  } finally {
    loading.value = false
  }
}

// Handle page change
const handlePageChange = (page, newPageSize) => {
  currentPage.value = page
  loadPosts(page, newPageSize)
  // Scroll to top when changing pages
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// Handle page size change
const handlePageSizeChange = (current, size) => {
  currentPage.value = 1 // Reset to first page when changing page size
  pageSize.value = size
  loadPosts(1, size)
}

// Search handler
const handleSearch = () => {
  // Search is reactive through computed property
  if (filteredPosts.value.length === 0) {
    message.info('No posts match your search')
  }
}

// Clear search and reload
watch(searchText, (newValue) => {
  if (newValue === '') {
    loadPosts()
  }
})

// View post details
const viewPost = (id) => {
  router.push(`/posts/${id}`)
}

// Edit post
const editPost = (id) => {
  router.push(`/posts/${id}/edit`)
}

// Delete post
const deletePost = async (id) => {
  try {
    await postService.deletePost(id)
    message.success('Post deleted successfully')
    
    // Check if we need to go to previous page
    if (posts.value.length === 1 && currentPage.value > 1) {
      currentPage.value -= 1
    }
    
    loadPosts(currentPage.value, pageSize.value)
  } catch (error) {
    message.error('Failed to delete post')
    console.error('Error deleting post:', error)
  }
}

// Truncate content for preview
const truncateContent = (content, maxLength = 150) => {
  if (!content) return ''
  if (content.length <= maxLength) return content
  return content.substring(0, maxLength) + '...'
}

// Format date
const formatDate = (date) => {
  return dayjs(date).format('MMM D, YYYY')
}

// Initialize
onMounted(() => {
  loadPosts()
})
</script>

<style scoped>
.post-list {
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

.page-header {
  background: #fff;
  margin-bottom: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.filter-card {
  margin-bottom: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.empty-state {
  background: #fff;
  padding: 80px 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.posts-grid {
  margin-top: 0;
}

.post-card {
  height: 100%;
  border-radius: 8px;
  transition: all 0.3s ease;
  border: 1px solid #f0f0f0;
}

.post-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transform: translateY(-4px);
}

.post-title {
  font-size: 18px;
  font-weight: 600;
  color: #1890ff;
  cursor: pointer;
  transition: color 0.3s;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
  min-height: 50px;
}

.post-title:hover {
  color: #40a9ff;
}

.post-content {
  color: rgba(0, 0, 0, 0.65);
  margin-top: 12px;
  margin-bottom: 16px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 72px;
}

.post-meta {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
  color: rgba(0, 0, 0, 0.45);
  font-size: 13px;
}

.post-meta :deep(.anticon) {
  margin-right: 6px;
}

.pagination-card {
  margin-top: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  text-align: center;
}

.custom-pagination {
  display: flex;
  justify-content: center;
}

.custom-pagination :deep(.ant-pagination-item) {
  border-radius: 4px;
}

.custom-pagination :deep(.ant-pagination-item-active) {
  border-color: #1890ff;
  background-color: #1890ff;
}

.custom-pagination :deep(.ant-pagination-item-active a) {
  color: #fff;
}
</style>
