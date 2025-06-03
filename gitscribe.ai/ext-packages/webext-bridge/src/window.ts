import { createEndpointRuntime } from './internal/endpoint-runtime'
import { usePostMessaging } from './internal/post-message'
import { createStreamWirings } from './internal/stream'

const win = usePostMessaging('window')

const endpointRuntime = createEndpointRuntime('window', message =>
  win.postMessage(message),
)

win.onMessage((msg) => {
  if ('type' in msg && 'transactionID' in msg)
    endpointRuntime.endTransaction(msg.transactionID)
  else endpointRuntime.handleMessage(msg)
})

export function setNamespace(nsps: string): void {
  win.setNamespace(nsps)
  win.enable()
}

export const { sendMessage, onMessage } = endpointRuntime
export const { openStream, onOpenStreamChannel }
  = createStreamWirings(endpointRuntime)
