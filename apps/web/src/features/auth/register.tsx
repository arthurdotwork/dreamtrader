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

const registrationSchema = z.object({
	email: z.string().email(),
	password: z.string().min(6),
	passwordConfirmation: z.string().min(6),
}).refine(data => data.password === data.passwordConfirmation, {
	message: 'Passwords do not match',
	path: ['passwordConfirmation']
});

const RegistrationScreen = () => {
	const router = useRouter();

	const form = useForm<z.infer<typeof registrationSchema>>({
		resolver: zodResolver(registrationSchema),
		defaultValues: {
			email: '',
			password: '',
			passwordConfirmation: ''
		}
	});

	const registerQuery = useMutation({
		mutationFn: ({ email, password }: {
			email: string,
			password: string
		}) => client.post('/api/v1/register', { email, password }, {}),
		onSuccess: (data) => {
			router.push(Routes.Auth);
		},
		onError: (error) => {
			alert(error);
		}
	});

	const onSubmit = (values: z.infer<typeof registrationSchema>) => {
		registerQuery.mutate(values);
	};

	return (
		<>
			<div className="border rounded-lg p-8 max-w-md w-full bg-white">
				<div className="mb-4">
					<h1 className="font-bold text-xl">Sign up</h1>
					<p className="text-sm text-muted-foreground">Create your DreamTrader account</p>
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
										This is the email address you will use to sign in.
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
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="passwordConfirmation"
							render={({ field }) => (
								<FormItem className="mb-4">
									<FormLabel>Repeat your password</FormLabel>
									<FormControl>
										<Input type="password"
													 placeholder="********" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<div className="mt-8 flex items-center space-x-4">
							<Button type="submit">Sign up</Button>
							<Link
								href="/auth"
								className="text-muted-foreground text-sm hover:text-black transition"
							>
								or sign in.
							</Link>
						</div>
					</form>
				</Form>
			</div>
		</>
	);
};

export default RegistrationScreen;