'use client';

import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
	Form,
	FormControl, FormDescription,
	FormField,
	FormItem,
	FormLabel, FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useMutation } from '@tanstack/react-query';
import { client } from '@/lib/client';
import Link from 'next/link';

const resetSchema = z.object({
	email: z.string().email(),
});

const ResetScreen = () => {
	const form = useForm<z.infer<typeof resetSchema>>({
		resolver: zodResolver(resetSchema),
		defaultValues: {
			email: '',
		},
	});

	const resetMutation = useMutation({
		mutationFn: ({ email }: {
			email: string
		}) => client.post('/api/v1/auth/reset', { email }, {}),
		onSuccess: (data) => {
			// TODO: Handle this endpoint.
		},
		onError: (error) => {
			alert(error);
		},
	});

	const onSubmit = (values: z.infer<typeof resetSchema>) => {
		resetMutation.mutate(values);
	};

	return (
		<>
			<div className="border rounded-lg p-8 max-w-md w-full bg-white">
				<div className="mb-4">
					<h1 className="font-bold text-xl">Reset your password</h1>
					<p className="text-sm text-muted-foreground">To gain back access your
						DreamTrader account</p>
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
						<div className="mt-8 flex items-center space-x-4">
							<Button type="submit">Reset my password</Button>
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

export default ResetScreen;