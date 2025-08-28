import toast from "react-hot-toast"

export type Notification = {
    type: "error" | "warning" | "info",
    message: string
}

export const notifyError = (message: string) => {
    notify({
        type: "error",
        message: message
    });
}

export const notifyWarning = (message: string) => {
    notify({
        type: "warning",
        message: message
    });
}

export const notifyInfo = (message: string) => {
    notify({
        type: "info",
        message: message
    });
}

export const notify = (notification: Notification) => {
    switch (notification.type) {
        case "error":
            toast.error(notification.message, {
                position: 'bottom-right'
            });
            break
        case "warning":
            toast(notification.message, {
                position: 'bottom-right',
            
            });
            break
        case "info":
            toast.success(notification.message, {
                position: 'bottom-right'
            });
            break
    }
}