"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

const FormSchema = z.object({
  systemName: z.string().min(1, {
    message: "Must be at least 1 character.",
  }),
});

export function SettingsForm() {
  // todo: get system name from API

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      systemName: "",
    },
  });

  function onSubmit(data: z.infer<typeof FormSchema>) {
    toast("You submitted the following values", {
      description: (
        <pre className="mt-2 w-[320px] rounded-md bg-neutral-950 p-4">
          <code className="text-white">{JSON.stringify(data, null, 2)}</code>
        </pre>
      ),
    });
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="my-4 w-full space-y-6"
      >
        <div className="max-w-[500px]">
          <FormField
            control={form.control}
            name="systemName"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-muted-foreground">
                  System Name
                </FormLabel>
                <FormControl>
                  <Input
                    className="max-w-[300px]"
                      placeholder="Enter a name for your system"
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  This is your system's name, which will be displayed in client
                  applications.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
}
