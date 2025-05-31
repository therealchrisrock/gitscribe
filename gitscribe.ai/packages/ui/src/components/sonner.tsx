"use client"

import type { ToasterProps as SonnerNativeToasterProps } from "sonner"
import React, { useEffect, useState, ComponentType } from "react"
import { useTheme } from "next-themes"

// Inner component that is only rendered on the client
const InnerClientToaster = ({ sonnerProps }: { sonnerProps: SonnerNativeToasterProps }) => {
  const { theme = "system" } = useTheme()
  const [ActualSonnerComponent, setActualSonnerComponent] = useState<ComponentType<SonnerNativeToasterProps> | null>(null)

  useEffect(() => {
    import("sonner").then((mod) => {
      // Ensure the imported component is correctly typed
      setActualSonnerComponent(() => mod.Toaster as ComponentType<SonnerNativeToasterProps>)
    })
  }, [])

  if (!ActualSonnerComponent) {
    return null // Or a loading placeholder
  }

  // Merge custom styles with any passed-in styles
  const combinedStyle = {
    "--normal-bg": "var(--popover)",
    "--normal-text": "var(--popover-foreground)",
    "--normal-border": "var(--border)",
    ...(sonnerProps.style || {}),
  } as React.CSSProperties

  return (
    <ActualSonnerComponent
      theme={theme as SonnerNativeToasterProps["theme"]}
      className="toaster group"
      {...sonnerProps} // Pass down all original props
      style={combinedStyle}
    />
  )
}

// Outer Toaster component that ensures client-side rendering for InnerClientToaster
const Toaster = (props: SonnerNativeToasterProps) => {
  const [isClient, setIsClient] = useState(false)

  useEffect(() => {
    setIsClient(true) // This effect runs only on the client
  }, [])

  if (!isClient) {
    // Render nothing on the server, or a non-interactive placeholder
    // For Toaster, which is a global notification system, null is usually fine.
    return null
  }

  // Pass all props to the inner client-side component
  return <InnerClientToaster sonnerProps={props} />
}

export { Toaster }
