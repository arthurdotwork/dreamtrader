'use client';

import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  Form,
  FormControl, FormDescription,
  FormField,
  FormItem,
  FormLabel, FormMessage
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useMutation } from '@tanstack/react-query';
import { client } from '@/lib/client';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import Routes from '@/routes';
import axios from 'axios';
import { AuthenticationResponse } from '@/types/auth';

const authenticationSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6)
});



const AuthenticationScreen = () => {
  const router = useRouter();

  const form = useForm<z.infer<typeof authenticationSchema>>({
    resolver: zodResolver(authenticationSchema),
    defaultValues: {
      email: '',
      password: ''
    }
  });

  const authenticateQuery = useMutation({
    mutationFn: ({ email, password }: {
      email: string,
      password: string
    }) => client.post<AuthenticationResponse>('/api/v1/auth', { email, password }, {}),
    onSuccess: async (data) => {
      await axios.post('/api/auth', {
        accessToken: data.accessToken.token,
        refreshToken: data.refreshToken.token
      });

      localStorage.setItem('accessToken', data.accessToken.token);
      localStorage.setItem('refreshToken', data.refreshToken.token);

      router.push(Routes.Dashboard);
    },
    onError: (error) => {
      alert(error);
    }
  });

  const onSubmit = (values: z.infer<typeof authenticationSchema>) => {
    authenticateQuery.mutate(values);
  };

  return (
    <>
      <div className="border rounded-lg p-8 max-w-md w-full bg-white">
        <div className="mb-4">
          <h1 className="font-bold text-xl">Sign-in</h1>
          <p className="text-sm text-muted-foreground">Access your DreamTrader account</p>
        </div>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem className="mb-4">
                  <FormLabel>Email address</FormLabel>
                  <FormControl>
                    <Input
                      type="email"
                      placeholder="mail@domain.tld" {...field}
                      autoComplete={'off'}
                    />
                  </FormControl>
                  <FormDescription>
                    This is the email address you used when
                    you signed up.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem className="mb-4">
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password"
                           placeholder="********" {...field} />
                  </FormControl>
                  <FormDescription>
                    <Link href={"/auth/reset"}
                       className="hover:text-black transition">Forgot
                      your password?
                    </Link>
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="mt-8 flex items-center space-x-4">
              <Button type="submit">Sign in</Button>
              <Link
                href="/register"
                className="text-muted-foreground text-sm hover:text-black transition"
              >
                or create an account.
              </Link>
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default AuthenticationScreen;