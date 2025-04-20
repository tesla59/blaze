export const getOrCreateClientId = () => {
    // Check if we're in a browser environment
    if (typeof window !== 'undefined') {
        let clientId = localStorage.getItem('clientId');
        if (!clientId) {
            clientId = crypto.randomUUID();
            localStorage.setItem('clientId', clientId);
        }
        return clientId;
    }
    // Fallback for server-side rendering
    return `client-${Date.now()}`;
};