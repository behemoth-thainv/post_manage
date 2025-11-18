<template>
  <div class="post-form">
    <!-- Page Header -->
    <a-page-header
      class="page-header"
      :title="isEditMode ? 'Edit Post' : 'Create New Post'"
      :sub-title="isEditMode ? 'Update your post' : 'Write a new blog post'"
      @back="$router.back()"
    />

    <!-- Loading State -->
    <div v-if="loading && isEditMode" class="loading-container">
      <a-spin size="large" tip="Loading post..." />
    </div>

    <!-- Form Card -->
    <a-card v-else class="form-card" :bordered="false">
      <a-form
        ref="formRef"
        :model="formState"
        :rules="rules"
        layout="vertical"
        @finish="handleSubmit"
      >
        <!-- Title Input -->
        <a-form-item label="Title" name="title" has-feedback>
          <a-input
            v-model:value="formState.title"
            placeholder="Enter post title"
            size="large"
            :maxlength="200"
            show-count
          >
            <template #prefix>
              <edit-outlined style="color: rgba(0, 0, 0, 0.25)" />
            </template>
          </a-input>
        </a-form-item>

        <!-- Content Textarea -->
        <a-form-item label="Content" name="content" has-feedback>
          <a-textarea
            v-model:value="formState.content"
            placeholder="Write your post content here..."
            :rows="12"
            :maxlength="5000"
            show-count
            style="resize: vertical"
          />
        </a-form-item>

        <!-- Preview Section -->
        <a-divider orientation="left">
          <eye-outlined /> Preview
        </a-divider>
        <a-card class="preview-card" :bordered="true">
          <template v-if="formState.title || formState.content">
            <h2 class="preview-title">{{ formState.title || 'Post Title' }}</h2>
            <div class="preview-content">
              {{ formState.content || 'Post content will appear here...' }}
            </div>
          </template>
          <a-empty v-else description="Start writing to see preview" :image="simpleImage" />
        </a-card>

        <!-- Form Actions -->
        <a-form-item class="form-actions">
          <a-space :size="16">
            <a-button
              type="primary"
              html-type="submit"
              size="large"
              :loading="submitting"
            >
              <template #icon>
                <save-outlined />
              </template>
              {{ isEditMode ? 'Update Post' : 'Create Post' }}
            </a-button>
            <a-button size="large" @click="handleReset">
              <template #icon>
                <redo-outlined />
              </template>
              Reset
            </a-button>
            <a-button size="large" @click="$router.back()">
              <template #icon>
                <close-outlined />
              </template>
              Cancel
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { Empty } from 'ant-design-vue'
import {
  EditOutlined,
  SaveOutlined,
  RedoOutlined,
  CloseOutlined,
  EyeOutlined,
} from '@ant-design/icons-vue'
import { postService } from '@/services/api'

const route = useRoute()
const router = useRouter()
const formRef = ref()
const loading = ref(false)
const submitting = ref(false)
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

// Check if edit mode
const isEditMode = computed(() => {
  return route.name === 'post-edit' && route.params.id
})

// Form state
const formState = reactive({
  title: '',
  content: '',
})

// Validation rules
const rules = {
  title: [
    {
      required: true,
      message: 'Please enter post title',
      trigger: 'blur',
    },
    {
      min: 3,
      message: 'Title must be at least 3 characters',
      trigger: 'blur',
    },
    {
      max: 200,
      message: 'Title must not exceed 200 characters',
      trigger: 'blur',
    },
  ],
  content: [
    {
      required: true,
      message: 'Please enter post content',
      trigger: 'blur',
    },
    {
      min: 10,
      message: 'Content must be at least 10 characters',
      trigger: 'blur',
    },
  ],
}

// Load post for editing
const loadPost = async () => {
  if (!isEditMode.value) return

  loading.value = true
  try {
    const response = await postService.getPost(route.params.id)
    formState.title = response.data.title
    formState.content = response.data.content
  } catch (error) {
    message.error('Failed to load post')
    console.error('Error loading post:', error)
    router.push('/')
  } finally {
    loading.value = false
  }
}

// Handle form submission
const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEditMode.value) {
      // Update existing post
      await postService.updatePost(route.params.id, formState)
      message.success('Post updated successfully')
    } else {
      // Create new post
      await postService.createPost(formState)
      message.success('Post created successfully')
    }
    router.push('/')
  } catch (error) {
    message.error(
      isEditMode.value ? 'Failed to update post' : 'Failed to create post'
    )
    console.error('Error submitting form:', error)
  } finally {
    submitting.value = false
  }
}

// Handle form reset
const handleReset = () => {
  if (isEditMode.value) {
    loadPost()
  } else {
    formRef.value.resetFields()
  }
  message.info('Form has been reset')
}

// Initialize
onMounted(() => {
  if (isEditMode.value) {
    loadPost()
  }
})
</script>

<style scoped>
.post-form {
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

.form-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.preview-card {
  background: #fafafa;
  border: 1px dashed #d9d9d9;
  margin-bottom: 24px;
}

.preview-title {
  font-size: 24px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 16px;
}

.preview-content {
  font-size: 15px;
  line-height: 1.8;
  color: rgba(0, 0, 0, 0.85);
  white-space: pre-wrap;
  word-wrap: break-word;
}

.form-actions {
  margin-top: 32px;
  margin-bottom: 0;
}

:deep(.ant-form-item-label > label) {
  font-weight: 600;
  font-size: 15px;
}

:deep(.ant-input),
:deep(.ant-input-textarea textarea) {
  font-size: 15px;
}
</style>
