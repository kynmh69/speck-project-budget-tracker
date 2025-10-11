import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { authApi, RegisterData, LoginData } from '@/lib/auth';
import { useAuthStore } from '@/store/auth-store';

export function useLogin() {
  const { setAuth } = useAuthStore();
  const router = useRouter();

  return useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      setAuth(data.user, data.token);
      router.push('/dashboard');
    },
  });
}

export function useRegister() {
  const { setAuth } = useAuthStore();
  const router = useRouter();

  return useMutation({
    mutationFn: authApi.register,
    onSuccess: (data) => {
      setAuth(data.user, data.token);
      router.push('/dashboard');
    },
  });
}

export function useLogout() {
  const { clearAuth } = useAuthStore();
  const router = useRouter();
  const queryClient = useQueryClient();

  return () => {
    clearAuth();
    queryClient.clear();
    router.push('/login');
  };
}

export function useCurrentUser() {
  const { isAuthenticated } = useAuthStore();

  return useQuery({
    queryKey: ['currentUser'],
    queryFn: authApi.me,
    enabled: isAuthenticated,
    retry: false,
  });
}
