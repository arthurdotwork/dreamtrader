import {createLazyFileRoute, Link} from '@tanstack/react-router'
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

const registerSchema = z.object({
    email: z.string().email(),
    password: z.string().min(8),
})

const Authenticate = () => {
    const registerForm = useForm<z.infer<typeof registerSchema>>({
        resolver: zodResolver(registerSchema),
        defaultValues: {
            email: "",
            password: "",
        },
    })

    const registerMutation = useMutation({
        mutationFn: (data: z.infer<typeof registerSchema>) => axios.post('http://localhost:8080/api/v1/register', data),
        onSuccess: () => {
            alert('success')
        }
    })

    const onSubmit = (data: z.infer<typeof registerSchema>) => {
        registerMutation.mutate(data)
    }

    return (
        <div className="p-6 max-w-2xl m-auto">
            <div className="mb-4">
                <h1 className="text-xl font-black">Connect to DreamTrader</h1>
                <p className="text-gray-600">Register to get access to all
                    DreamTrader features</p>
            </div>
            <Form {...registerForm}>
                <form onSubmit={registerForm.handleSubmit(onSubmit)}
                      className="space-y-8">
                    <FormField
                        control={registerForm.control}
                        name="email"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Email</FormLabel>
                                <FormControl>
                                    <Input placeholder="mail@domain.tld" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={registerForm.control}
                        name="password"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Password</FormLabel>
                                <FormControl>
                                    <Input type="password" placeholder="********" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <div className="flex items-center">
                        <Button type="submit">Register</Button>
                        <p className="text-gray-600 ml-2">Don't have an account? <Link className="underline" to={'/register'}>Sign
                            up</Link> instead.</p>
                    </div>
                </form>
            </Form>
        </div>
    )
}

export const Route = createLazyFileRoute('/authenticate')({
    component: Authenticate,
})

