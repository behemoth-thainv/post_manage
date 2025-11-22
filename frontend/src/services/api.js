import axios from 'axios'

const baseURL = import.meta?.env?.VITE_API_BASE || 'http://localhost:3000/api'

const apiClient = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const postService = {
  /**
   * Get all posts with pagination
   * @param {number} page - Page number (default: 1)
   * @param {number} limit - Items per page (default: 12)
   */
  getPosts(page = 1, limit = 12) {
    return apiClient.get('/posts', {
      params: { page, limit }
    })
  },

  /**
   * Get a single post by ID
   */
  getPost(id) {
    return apiClient.get(`/posts/${id}`)
  },

  /**
   * Create a new post
   */
  createPost(postData) {
    return apiClient.post('/posts', postData)
  },

  /**
   * Update an existing post
   */
  updatePost(id, postData) {
    return apiClient.put(`/posts/${id}`, postData)
  },

  /**
   * Delete a post
   */
  deletePost(id) {
    return apiClient.delete(`/posts/${id}`)
  },
}

export default apiClient
