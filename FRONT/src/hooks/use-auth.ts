import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { loginUser, registerUser } from '@/services/auth'

export function useLogin() {
  const navigate = useNavigate()

  return useMutation({
    mutationFn: loginUser,
    onSuccess: (data) => {
      console.log('Login realizado com sucesso:', data)

      if (data.token) {
        localStorage.setItem('authToken', data.token)
      }

      toast.success('Login realizado com sucesso!', {
        description: 'Bem-vindo!'
      })

      navigate('/dashboard')
    },
    onError: (error) => {
      console.error('Erro no login:', error)
      toast.error('Erro no login', {
        description: error.message
      })
    }
  })
}

export function useRegister() {
  const navigate = useNavigate()

  return useMutation({
    mutationFn: registerUser,
    onSuccess: (data) => {
      console.log('Registro realizado com sucesso:', data)

      toast.success('Conta criada com sucesso!', {
        description: `Bem-vindo, ${data.email}! Faça login para continuar.`
      })

      navigate('/login')
    },
    onError: (error) => {
      console.error('Erro no registro:', error)
      toast.error('Erro no registro', {
        description: error.message
      })
    }
  })
}
