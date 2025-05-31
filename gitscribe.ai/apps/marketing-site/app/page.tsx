import Link from "next/link"
import { ArrowRight, Check, Github, Headphones, MessageSquare, Ticket } from "lucide-react"
import { Button } from "@workspace/ui/components/button"

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen">
      <header className="px-4 lg:px-6 h-16 flex items-center border-b">
        <Link className="flex items-center justify-center" href="#">
          <span className="font-bold text-xl flex items-center gap-2">
            <Github className="h-6 w-6" />
            GitScribe
          </span>
        </Link>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="#features">
            Features
          </Link>
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="#how-it-works">
            How It Works
          </Link>
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="#pricing">
            Pricing
          </Link>
        </nav>
      </header>
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48 bg-gradient-to-b from-purple-50 to-white dark:from-gray-900 dark:to-gray-950">
          <div className="container px-4 md:px-6">
            <div className="grid gap-6 lg:grid-cols-[1fr_400px] lg:gap-12 xl:grid-cols-[1fr_600px]">
              <div className="flex flex-col justify-center space-y-4">
                <div className="space-y-2">
                  <h1 className="text-3xl font-bold tracking-tighter sm:text-5xl xl:text-6xl/none">
                    Transform Technical Calls into GitHub Tickets
                  </h1>
                  <p className="max-w-[600px] text-gray-500 md:text-xl dark:text-gray-400">
                    GitScribe records, transcribes, and summarizes your engineering meetings, then automatically
                    generates actionable GitHub tickets.
                  </p>
                </div>
                <div className="flex flex-col gap-2 min-[400px]:flex-row">
                  <Button className="bg-purple-600 hover:bg-purple-700">
                    Install Extension <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                  <Button variant="outline">Watch Demo</Button>
                </div>
              </div>
              <div className="flex items-center justify-center">
                <div className="relative w-full max-w-[500px] aspect-video rounded-xl border bg-white shadow-xl dark:bg-gray-800">
                  <div className="absolute top-0 left-0 right-0 h-8 bg-gray-100 dark:bg-gray-700 rounded-t-xl flex items-center px-4">
                    <div className="flex space-x-2">
                      <div className="w-3 h-3 rounded-full bg-red-500"></div>
                      <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                      <div className="w-3 h-3 rounded-full bg-green-500"></div>
                    </div>
                    <div className="mx-auto text-xs text-gray-500 dark:text-gray-400">
                      Google Meet - Engineering Call
                    </div>
                  </div>
                  <div className="pt-10 p-4 flex flex-col h-full">
                    <div className="flex-1 flex items-center justify-center">
                      <div className="text-center space-y-2">
                        <div className="inline-flex items-center justify-center p-3 bg-purple-100 rounded-full dark:bg-purple-900">
                          <Headphones className="h-6 w-6 text-purple-600 dark:text-purple-300" />
                        </div>
                        <p className="text-sm text-gray-500 dark:text-gray-400">Recording in progress...</p>
                        <div className="flex justify-center space-x-1">
                          {[1, 2, 3, 4, 5].map((i) => (
                            <div
                              key={i}
                              className="w-1 h-3 bg-purple-600 dark:bg-purple-400 rounded-full animate-pulse"
                              style={{ animationDelay: `${i * 0.1}s` }}
                            ></div>
                          ))}
                        </div>
                      </div>
                    </div>
                    <div className="mt-4 p-3 bg-gray-100 dark:bg-gray-700 rounded-lg text-xs">
                      <div className="font-medium">GitScribe is active</div>
                      <div className="text-gray-500 dark:text-gray-400">Transcribing and analyzing your meeting</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section id="features" className="w-full py-12 md:py-24 lg:py-32 bg-white dark:bg-gray-950">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <div className="space-y-2">
                <div className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-purple-600 text-white">
                  Features
                </div>
                <h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">
                  Supercharge Your Engineering Workflow
                </h2>
                <p className="max-w-[900px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed dark:text-gray-400">
                  GitScribe seamlessly integrates with your browser to transform technical discussions into actionable
                  tasks.
                </p>
              </div>
            </div>
            <div className="mx-auto grid max-w-5xl items-center gap-6 py-12 lg:grid-cols-3 lg:gap-12">
              <div className="flex flex-col justify-center space-y-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-purple-100 dark:bg-purple-900">
                  <Headphones className="h-6 w-6 text-purple-600 dark:text-purple-300" />
                </div>
                <div className="space-y-2">
                  <h3 className="text-xl font-bold">Smart Recording</h3>
                  <p className="text-gray-500 dark:text-gray-400">
                    Automatically records your Google Meet engineering calls with a single click.
                  </p>
                </div>
              </div>
              <div className="flex flex-col justify-center space-y-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-purple-100 dark:bg-purple-900">
                  <MessageSquare className="h-6 w-6 text-purple-600 dark:text-purple-300" />
                </div>
                <div className="space-y-2">
                  <h3 className="text-xl font-bold">AI Transcription</h3>
                  <p className="text-gray-500 dark:text-gray-400">
                    Converts spoken discussions into accurate text with speaker identification.
                  </p>
                </div>
              </div>
              <div className="flex flex-col justify-center space-y-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-purple-100 dark:bg-purple-900">
                  <Ticket className="h-6 w-6 text-purple-600 dark:text-purple-300" />
                </div>
                <div className="space-y-2">
                  <h3 className="text-xl font-bold">GitHub Integration</h3>
                  <p className="text-gray-500 dark:text-gray-400">
                    Automatically generates well-structured tickets based on meeting content.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section id="how-it-works" className="w-full py-12 md:py-24 lg:py-32 bg-gray-50 dark:bg-gray-900">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <div className="space-y-2">
                <div className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-purple-600 text-white">
                  How It Works
                </div>
                <h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">Simple, Powerful Workflow</h2>
                <p className="max-w-[900px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed dark:text-gray-400">
                  GitScribe works seamlessly in the background, turning your discussions into actionable tasks.
                </p>
              </div>
            </div>
            <div className="mx-auto grid max-w-5xl gap-8 py-12 lg:grid-cols-4">
              {[
                {
                  step: "1",
                  title: "Install & Connect",
                  description: "Add GitScribe to your browser and connect your GitHub account.",
                  icon: <Github className="h-6 w-6" />,
                },
                {
                  step: "2",
                  title: "Join Your Meeting",
                  description: "Start or join a Google Meet call and activate GitScribe.",
                  icon: <Headphones className="h-6 w-6" />,
                },
                {
                  step: "3",
                  title: "AI Transcription",
                  description: "GitScribe listens and transcribes the entire conversation.",
                  icon: <MessageSquare className="h-6 w-6" />,
                },
                {
                  step: "4",
                  title: "Generate Tickets",
                  description: "Review and approve AI-generated GitHub tickets.",
                  icon: <Ticket className="h-6 w-6" />,
                },
              ].map((item, index) => (
                <div key={index} className="flex flex-col items-center space-y-4 relative">
                  <div className="flex h-16 w-16 items-center justify-center rounded-full bg-purple-600 text-white">
                    {item.icon}
                  </div>
                  <div className="absolute -top-3 -right-3 flex h-8 w-8 items-center justify-center rounded-full bg-gray-200 text-gray-800 font-bold dark:bg-gray-700 dark:text-gray-200">
                    {item.step}
                  </div>
                  <div className="space-y-2 text-center">
                    <h3 className="text-xl font-bold">{item.title}</h3>
                    <p className="text-gray-500 dark:text-gray-400">{item.description}</p>
                  </div>
                  {index < 3 && (
                    <div className="hidden lg:block absolute top-8 left-full w-full h-0.5 bg-gray-200 dark:bg-gray-700 -z-10">
                      <ArrowRight className="absolute right-0 top-1/2 -translate-y-1/2 text-gray-300 dark:text-gray-600" />
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </section>

        <section id="pricing" className="w-full py-12 md:py-24 lg:py-32 bg-white dark:bg-gray-950">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <div className="space-y-2">
                <div className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-purple-600 text-white">
                  Pricing
                </div>
                <h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">Simple, Transparent Pricing</h2>
                <p className="max-w-[900px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed dark:text-gray-400">
                  Choose the plan that works best for your team.
                </p>
              </div>
            </div>
            <div className="mx-auto grid max-w-5xl gap-6 py-12 lg:grid-cols-3">
              {[
                {
                  name: "Free",
                  price: "$0",
                  description: "Perfect for trying out GitScribe.",
                  features: [
                    "5 recordings per month",
                    "Basic transcription",
                    "Manual ticket creation",
                    "1 GitHub repository",
                  ],
                  cta: "Get Started",
                  highlighted: false,
                },
                {
                  name: "Pro",
                  price: "$19",
                  period: "per user/month",
                  description: "Everything you need for a small team.",
                  features: [
                    "Unlimited recordings",
                    "Advanced AI transcription",
                    "Automatic ticket generation",
                    "Multiple GitHub repositories",
                    "Priority support",
                  ],
                  cta: "Start Free Trial",
                  highlighted: true,
                },
                {
                  name: "Enterprise",
                  price: "Custom",
                  description: "For organizations with advanced needs.",
                  features: [
                    "Everything in Pro",
                    "Custom integrations",
                    "Advanced security",
                    "Dedicated account manager",
                    "SLA guarantees",
                  ],
                  cta: "Contact Sales",
                  highlighted: false,
                },
              ].map((plan, index) => (
                <div
                  key={index}
                  className={`flex flex-col p-6 ${plan.highlighted
                    ? "rounded-xl border-2 border-purple-600 bg-white shadow-lg dark:bg-gray-800"
                    : "rounded-xl border bg-white shadow-md dark:bg-gray-800"
                    }`}
                >
                  <div className="space-y-2">
                    <h3 className="text-2xl font-bold">{plan.name}</h3>
                    <div className="space-y-1">
                      <div className="text-3xl font-bold">{plan.price}</div>
                      {plan.period && <p className="text-sm text-gray-500 dark:text-gray-400">{plan.period}</p>}
                    </div>
                    <p className="text-gray-500 dark:text-gray-400">{plan.description}</p>
                  </div>
                  <ul className="my-6 space-y-2 flex-1">
                    {plan.features.map((feature, featureIndex) => (
                      <li key={featureIndex} className="flex items-center">
                        <Check className="mr-2 h-4 w-4 text-purple-600" />
                        <span>{feature}</span>
                      </li>
                    ))}
                  </ul>
                  <Button
                    className={plan.highlighted ? "bg-purple-600 hover:bg-purple-700" : "bg-gray-900 dark:bg-gray-700"}
                  >
                    {plan.cta}
                  </Button>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section className="w-full py-12 md:py-24 lg:py-32 bg-purple-600 text-white">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <div className="space-y-2">
                <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
                  Ready to Transform Your Workflow?
                </h2>
                <p className="mx-auto max-w-[700px] text-white/80 md:text-xl">
                  Join thousands of engineering teams who are saving time and improving collaboration with GitScribe.
                </p>
              </div>
              <div className="flex flex-col gap-2 min-[400px]:flex-row">
                <Button className="bg-white text-purple-600 hover:bg-gray-100">
                  Install GitScribe <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
                <Button variant="outline" className="border-white text-white hover:bg-purple-700">
                  Schedule Demo
                </Button>
              </div>
            </div>
          </div>
        </section>
      </main>
      <footer className="flex flex-col gap-2 sm:flex-row py-6 w-full border-t px-4 md:px-6">
        <p className="text-xs text-gray-500 dark:text-gray-400">© 2023 GitScribe. All rights reserved.</p>
        <nav className="sm:ml-auto flex gap-4 sm:gap-6">
          <Link className="text-xs hover:underline underline-offset-4 text-gray-500 dark:text-gray-400" href="#">
            Terms of Service
          </Link>
          <Link className="text-xs hover:underline underline-offset-4 text-gray-500 dark:text-gray-400" href="#">
            Privacy
          </Link>
          <Link className="text-xs hover:underline underline-offset-4 text-gray-500 dark:text-gray-400" href="#">
            Contact
          </Link>
        </nav>
      </footer>
    </div>
  )
}
