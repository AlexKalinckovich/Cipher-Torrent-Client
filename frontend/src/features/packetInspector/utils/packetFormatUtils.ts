export const formatTime = (isoString: string): string => {
    return new Date(isoString).toISOString().substring(11, 23);
};