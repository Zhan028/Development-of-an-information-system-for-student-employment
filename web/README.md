# Student Employment Web Application

Frontend application for the Student Employment Information System.

## Tech Stack

- **React 19** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework

## Project Structure

```
src/
├── api/          # API integration
├── components/   # Reusable UI components
├── pages/        # Page components
├── store/        # State management
├── hooks/        # Custom React hooks
├── types/        # TypeScript types
├── utils/        # Utility functions
├── routes/       # Routing config
├── App.tsx       # Root component
├── main.tsx      # Entry point
└── index.css     # Global styles + Tailwind
```

## Getting Started

### Prerequisites

- Node.js 20+ 
- npm or yarn

### Installation

```bash
npm install
```

### Development

```bash
npm run dev
```

The application will be available at `http://localhost:5173`

### Build

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

- `VITE_API_URL` - Backend API URL
- `VITE_API_GATEWAY_URL` - API Gateway URL
- `VITE_ENV` - Environment (development/production)

## Docker

Build and run with Docker:

```bash
docker build -t student-employment-web .
docker run -p 80:80 student-employment-web
```
