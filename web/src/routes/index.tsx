import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import App from '../App';
import { HomePage, LoginPage, RegisterPage } from '../pages';
import { AuthProvider } from '../context';

const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/register',
    element: <RegisterPage />,
  },
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <HomePage /> },
    ],
  },
]);

export const AppRouter = () => (
  <AuthProvider>
    <RouterProvider router={router} />
  </AuthProvider>
);
