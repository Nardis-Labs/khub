import { useEffect, useState } from 'react';
import { PieChartData } from './app/components/PieChart/PieChart';

export const getWindowDimensions = () => {
  const { innerWidth: width, innerHeight: height } = window;
  return {
    width,
    height
  };
};

export const useWindowDimensions = () => {
  const [windowDimensions, setWindowDimensions] = useState(getWindowDimensions());

  useEffect(() => {
    function handleResize() {
      setWindowDimensions(getWindowDimensions());
    }

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return windowDimensions;
};

export const renderPodStats = (data: any, filterFunc: (filter: string) => void): PieChartData[] => {
  const statusMap = new Map<string, number>();
  const statuses: PieChartData[] = [];
  if (data.length === 0){
    statuses.push({name: 'nil', value: 1, fill: 'gray', filterFunc: filterFunc, filterString: ''});
    return statuses;
  }

  for (let i = 0; i < data.length; i++){
    if (data[i].data.metadata.deletionTimestamp !== undefined) {
      statusMap.set('Terminating', (statusMap.get('Terminating') || 0) + 1);
    } else {
      statusMap.set(data[i].data.status.phase, (statusMap.get(data[i].data.status.phase) || 0) + 1);
    }
  }

  statusMap.forEach((value, key) => {
    let fill = '';
    let filterString = '';
    switch (key) {
      case 'Running':
        fill = '#00bc62';
        filterString = 'running';
        break;
      case 'Terminating':
        fill = '#00bc62';
        filterString = 'terminating';
        break;
      case 'Pending':
        fill = 'orange';
        filterString = 'pending';
        break;
      case 'Succeeded':
        fill = '#0088FE';
        filterString = 'succeeded';
        break;
      case 'Failed':
        fill = 'red';
        filterString = 'failed';
        break;
    }
    statuses.push({name: key, value: value, fill: fill, filterFunc: filterFunc, filterString: filterString});
  });
  return statuses;
};

export const renderDeploymentStats = (data: any, filterFunc: (filter: string) => void, kind: string): PieChartData[] => {
  const statusMap = new Map<string, number>();
  const statuses: PieChartData[] = [];
  if (data.length === 0){
    statuses.push({name: 'nil', value: 1, fill: 'gray', filterFunc: filterFunc, filterString: ''});
    return statuses;
  }

  for (let i = 0; i < data.length; i++){
    if (data[i].data.status.availableReplicas === data[i].data.status.replicas){
      statusMap.set('Running', (statusMap.get('Running') || 0) + 1);
    } else {
      statusMap.set('Pending', (statusMap.get('Pending') || 0) + 1);
    }
  }

  if (kind === "replicaset") {
    statusMap.forEach((value, key) => {
      let fill = '';
      let keyOverride = '';
      switch (key) {
        case 'Running':
          fill = '#00bc62';
          keyOverride = 'Active';
          break;
        case 'Pending':
          fill = 'grey';
          keyOverride = 'Inactive';
          break;
      }
      statuses.push({name: keyOverride, value: value, fill: fill, filterFunc: filterFunc, filterString: ''});
    });
  } else {
    statusMap.forEach((value, key) => {
      let fill = '';
      let filterString = '';
      switch (key) {
        case 'Running':
          fill = '#00bc62';
          filterString = 'running';
          break;
        case 'Pending':
          fill = 'orange';
          filterString = 'pending';
          break;
      }
      statuses.push({name: key, value: value, fill: fill, filterFunc: filterFunc, filterString: filterString});
    });
  }
  
  return statuses;
};

export const renderJobStatus = (data: any, filterFunc: (filter: string) => void): PieChartData[] => {
  const statusMap = new Map<string, number>();
  const statuses: PieChartData[] = [];
  if (data.length === 0){
    statuses.push({name: 'nil', value: 1, fill: 'gray', filterFunc: filterFunc, filterString: ''});
    return statuses;
  }

  for (let i = 0; i < data.length; i++){
    if (data[i].data.status.succeeded >= 1){
      statusMap.set('Succeeded', (statusMap.get('Succeeded') || 0) + 1);
    } else if (data[i].data.status.terminating >= 1 || data[i].data.status.ready >= 1 || data[i].data.status.active >= 1){
      statusMap.set('Pending', (statusMap.get('Pending') || 0) + 1);
    } else if (data[i].data.status.failed >= 1){
      statusMap.set('Failed', (statusMap.get('Failed') || 0) + 1);
    } 
  }
  statusMap.forEach((value, key) => {
    let fill = '';
    switch (key) {
      case 'Succeeded':
        fill = '#00bc62';
        break;
      case 'Pending':
        fill = 'orange';
        break;
      case 'Failed':
        fill = 'red';
        break;
    }
    statuses.push({name: key, value: value, fill: fill, filterFunc: filterFunc, filterString: ''});
  });
  return statuses;
};

export const renderCronJobStatus = (data: any, filterFunc: (filter: string) => void): PieChartData[] => {
  const statusMap = new Map<string, number>();
  const statuses: PieChartData[] = [];
  if (data.length === 0){
    statuses.push({name: 'nil', value: 1, fill: 'gray', filterFunc: filterFunc, filterString: ''});
    return statuses;
  }

  for (let i = 0; i < data.length; i++){
    if (data[i].data.spec.suspend !== true){
      statusMap.set('Active', (statusMap.get('Active') || 0) + 1);
    } else if (data[i].data.spec.suspend === false){
      statusMap.set('Suspended', (statusMap.get('Suspended') || 0) + 1);
    }
  }
  statusMap.forEach((value, key) => {
    let fill = '';
    switch (key) {
      case 'Active':
        fill = '#00bc62';
        break;
      case 'Suspended':
        fill = 'orange';
        break;
    }
    statuses.push({name: key, value: value, fill: fill, filterFunc: filterFunc, filterString: ''});
  });
  return statuses;
};