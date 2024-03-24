import {createLazyFileRoute, Link, useNavigate} from '@tanstack/react-router'
import {useMutation} from "@tanstack/react-query";
import {z} from "zod";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";

import { Button } from "@/components/ui/button"
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import axios from "axios";
import {AuthAccessToken, AuthRefreshToken, useAuth} from "@/contexts/auth.tsx";

const authenticateSchema = z.object({
    email: z.string().email(),
    password: z.string().min(8),
})

type AuthenticationResponse = {
    access_token: AuthAccessToken;
    refresh_token: AuthRefreshToken;
}

const Authenticate = () => {
    const navigate = useNavigate();
    const { saveCredentials, isAuthenticated } = useAuth();

    if (isAuthenticated()) {
        navigate({to: '/assets'})
    }

    const authenticateForm = useForm<z.infer<typeof authenticateSchema>>({
        resolver: zodResolver(authenticateSchema),
        defaultValues: {
            email: "",
            password: "",
        },
    })

    const authenticateMutation = useMutation({
        mutationFn: (data: z.infer<typeof authenticateSchema>) => axios.post<AuthenticationResponse>('http://localhost:8080/api/v1/authenticate', data),
        onSuccess: ({data}) => {
            saveCredentials({accessToken: data.access_token, refreshToken: data.refresh_token})
            navigate({to: '/assets'})
        }
    })

    const onSubmit = (data: z.infer<typeof authenticateSchema>) => {
        authenticateMutation.mutate(data)
    }

    return (
        <div className="h-screen flex items-center">
            <div className="p-6 max-w-2xl m-auto border rounded-lg">
                <div className="mb-4">
                    <h1 className="text-xl font-black">Connect to
                        DreamTrader</h1>
                    <p className="text-gray-600">Register to get access to all
                        DreamTrader features</p>
                </div>
                <Form {...authenticateForm}>
                    <form onSubmit={authenticateForm.handleSubmit(onSubmit)}
                          className="space-y-8">
                        <FormField
                            control={authenticateForm.control}
                            name="email"
                            render={({field}) => (
                                <FormItem>
                                    <FormLabel>Email</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder="mail@domain.tld" {...field} />
                                    </FormControl>
                                    <FormMessage/>
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={authenticateForm.control}
                            name="password"
                            render={({field}) => (
                                <FormItem>
                                    <FormLabel>Password</FormLabel>
                                    <FormControl>
                                        <Input type="password"
                                               placeholder="********" {...field} />
                                    </FormControl>
                                    <FormMessage/>
                                </FormItem>
                            )}
                        />

                        <div className="flex items-center">
                            <Button type="submit">Register</Button>
                            <p className="text-gray-600 ml-2">Don't have an
                                account? <Link className="underline"
                                               to={'/register'}>Sign
                                    up</Link> instead.</p>
                        </div>
                    </form>
                </Form>
            </div>
        </div>
    )
}

export const Route = createLazyFileRoute('/authenticate')({
    component: Authenticate,
})

