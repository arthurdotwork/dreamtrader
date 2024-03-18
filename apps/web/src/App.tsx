import {Button} from "@/components/ui/button.tsx";
import {z} from "zod";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {
    Form,
    FormControl,
    FormField,
    FormItem, FormLabel,
    FormMessage
} from "@/components/ui/form.tsx";
import {Input} from "@/components/ui/input.tsx";
import dreamTraderLogo from "@/assets/dreamtrader.svg";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";

const signInSchema = z.object({
    email: z.string().email(),
    password: z.string().min(8),
});
const App = () => {
    const form = useForm<z.infer<typeof signInSchema>>({
        resolver: zodResolver(signInSchema),
        defaultValues: {
            email: "",
            password: "",
        },
    })

    const onSubmit = (values: z.infer<typeof signInSchema>) => {
        console.log({values})
    }


    return (
        <div className="flex justify-center bg-zinc-50">
            <div className="w-1/3 h-screen flex flex-col justify-center">
                <div className="flex justify-center">
                    <img src={dreamTraderLogo} alt="DreamTrader" className="w-10 h-10 mr-4"/>
                    <h1 className="text-4xl text-black font-black">DreamTrader</h1>
                </div>
                <p className="text-muted-foreground mt-2 text-center">Improve your trading skills with fictional money</p>

                <Card className="w-full p-4 mt-8">
                    <CardContent className="grid gap-4 p-0">
                        <Form {...form}>
                            <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
                                <FormField
                                    control={form.control}
                                    name="email"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Email</FormLabel>
                                            <FormControl>
                                                <Input placeholder="user@gmail.com" {...field} />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="password"
                                    render={({field}) => (
                                        <FormItem className="mt-4">
                                            <FormLabel>Password</FormLabel>
                                            <FormControl>
                                                <Input type="password" placeholder="********" {...field} />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                                <Button type="submit" className="mt-4 bg-blue-800 hover:bg-blue-800">Access DreamTrader</Button>
                            </form>
                        </Form>
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}

export default App
