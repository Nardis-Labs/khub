export const wsConnect = async (resourceApi: string, cacheDataLoaded: any, updateCachedData: any, cacheEntryRemoved: any, pingInterval = 5000) => {
  // create a websocket connection when the cache subscription starts
  const ws = new WebSocket(resourceApi);
  
  setInterval(() => {ws.send('ping');}, pingInterval);
  try {
    // wait for the initial query to resolve before proceeding
    await cacheDataLoaded;

    // when data is received from the socket connection to the server,
    // if it is a message and for the appropriate channel,
    // update our query result with the received message
    const listener = (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      updateCachedData(() => {
        return data;
      });
    };

    ws.addEventListener('message', listener);
    
  } catch {
    // no-op in case `cacheEntryRemoved` resolves before `cacheDataLoaded`,
    // in which case `cacheDataLoaded` will throw
  }
  // cacheEntryRemoved will resolve when the cache subscription is no longer active
  await cacheEntryRemoved;
  // perform cleanup steps once the `cacheEntryRemoved` promise resolves
  ws.close();
};