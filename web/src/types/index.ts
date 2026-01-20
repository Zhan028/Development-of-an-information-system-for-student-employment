// TypeScript types
// Define shared types and interfaces here

// Example: User types
export interface User {
  id: string
  email: string
  name: string
  role: 'student' | 'employer' | 'admin'
}

// Example: Vacancy types
export interface Vacancy {
  id: string
  title: string
  description: string
  employerId: string
  createdAt: string
  updatedAt: string
}

// Example: Application types
export interface Application {
  id: string
  studentId: string
  vacancyId: string
  status: 'pending' | 'accepted' | 'rejected'
  createdAt: string
}
