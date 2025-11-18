# Post Management Frontend

A modern, responsive Vue.js application for managing blog posts with a beautiful UI powered by Ant Design Vue.

## Features

- 📝 **Full CRUD Operations**: Create, Read, Update, and Delete posts
- 🎨 **Beautiful UI**: Modern design with Ant Design Vue components
- 🔍 **Search Functionality**: Real-time search through posts
- 📱 **Responsive Design**: Works perfectly on all devices
- ⚡ **Fast Performance**: Built with Vite for lightning-fast development
- 🎯 **Clean Code**: Well-organized, maintainable code structure
- 🧭 **Vue Router**: Clear routing structure for easy feature additions

## Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **Ant Design Vue 4** - Enterprise-class UI components
- **Vue Router 4** - Official router for Vue.js
- **Axios** - HTTP client for API requests
- **Vite** - Next generation frontend tooling
- **Day.js** - Fast date manipulation library

## Project Structure

```
frontend/
├── src/
│   ├── layouts/          # Layout components
│   │   └── MainLayout.vue
│   ├── views/            # Page components
│   │   ├── PostList.vue
│   │   ├── PostDetail.vue
│   │   └── PostForm.vue
│   ├── services/         # API service layer
│   │   └── api.js
│   ├── router/           # Vue Router configuration
│   │   └── index.js
│   ├── styles/           # Global styles
│   │   └── main.css
│   ├── App.vue          # Root component
│   └── main.js          # Application entry point
├── index.html
├── vite.config.js
└── package.json
```

## Prerequisites

- Node.js 16+ and npm/yarn
- Backend API running on port 3000

## Installation

1. Install dependencies:
```bash
npm install
```

2. Start development server:
```bash
npm run dev
```

3. Open your browser and navigate to:
```
http://localhost:5173
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## API Integration

The application connects to a backend API on port 3000. API endpoints:

- `GET /posts` - Get all posts
- `GET /posts/:id` - Get single post
- `POST /posts` - Create new post
- `PUT /posts/:id` - Update post
- `DELETE /posts/:id` - Delete post

## Features Overview

### Post List
- Grid view of all posts
- Real-time search functionality
- Quick actions (view, edit, delete)
- Responsive card layout

### Post Detail
- Full post content display
- Creation and update timestamps
- Quick edit and delete actions

### Post Form
- Create new posts
- Edit existing posts
- Form validation
- Live preview
- Auto-save draft (can be added)

## Code Quality

- Clean, readable code with proper comments
- Consistent naming conventions
- Modular component structure
- Reusable service layer
- Proper error handling
- Loading states for better UX

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Future Enhancements

- Rich text editor
- Image upload
- Categories and tags
- User authentication
- Comments system
- Dark mode
- Pagination
- Advanced filters

## Contributing

Feel free to add new features or improve existing ones. The code structure is designed to be easily extensible.

## License

MIT
