# Prompt Explorer Frontend

A React-based frontend for the Prompt Explorer application, built with TypeScript, Vite, and Tailwind CSS.

## 🚀 Features

- **Modern React 18** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS** for styling
- **Axios** for API communication
- **Lucide React** for icons
- **Responsive design** for all devices

## 📁 Project Structure

```
fe/
├── src/
│   ├── components/          # React components
│   │   ├── SystemPromptEditor.tsx
│   │   ├── ChatInterface.tsx
│   │   ├── MessageBubble.tsx
│   │   ├── NavigationControls.tsx
│   │   └── LoadingIndicator.tsx
│   ├── types/               # TypeScript type definitions
│   ├── utils/               # Utility functions
│   │   └── api.ts          # API client
│   ├── hooks/              # Custom React hooks
│   ├── App.tsx             # Main application component
│   ├── main.tsx            # Application entry point
│   └── index.css           # Global styles
├── public/                 # Static assets
├── package.json            # Dependencies and scripts
├── vite.config.ts          # Vite configuration
├── tailwind.config.js      # Tailwind CSS configuration
└── tsconfig.json           # TypeScript configuration
```

## 🛠️ Development

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The development server will start on `http://localhost:3000` and proxy API requests to the Go backend at `http://localhost:8081`.

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## 🎨 Design System

### Components

- **SystemPromptEditor**: Large textarea for system prompt editing
- **ChatInterface**: Chat messages display and input
- **MessageBubble**: Individual message display
- **NavigationControls**: NEW/PREV/NEXT buttons
- **LoadingIndicator**: Loading state display

### Styling

- **Tailwind CSS** for utility-first styling
- **Custom components** with consistent design
- **Responsive design** for mobile and desktop
- **Smooth animations** and transitions

## 🔌 API Integration

The frontend communicates with the Go backend through:

- **GET /api/conversations** - List conversations
- **GET /api/conversations/{id}** - Get conversation details
- **POST /api/conversations** - Create new conversation
- **POST /api/conversations/{id}/messages** - Send message
- **POST /api/conversations/{id}/favorite** - Toggle favorite

## 📱 Responsive Design

- **Desktop**: Two-panel layout (system prompt + chat)
- **Mobile**: Stacked vertical layout
- **Touch-friendly** buttons and inputs
- **Optimized spacing** for all screen sizes

## 🚀 Deployment

Build the production bundle:

```bash
npm run build
```

The built files will be in the `dist/` directory, ready for deployment.

