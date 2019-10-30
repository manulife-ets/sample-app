import React, { useState } from 'react';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import { ThemeProvider } from 'styled-components';
import Grid from '@radd/react/dist/components/containers/Grid';
import { manulifeTheme } from '@radd/react/dist/components/base/Themes';
import ShortenerForm from './components/ShortenerForm';


const App = () => {
  const indexMap = new Map([
    [0, 'SHORTEN'],
    [1, 'CREATE'],
    [2, 'UPDATE']
  ]);

  const [task, setTask] = useState(indexMap.get(0))
  
  return (
  <ThemeProvider theme={manulifeTheme}>
    <Grid sm="1fr" padding="lg">
      <h1>Go.Manulife <em>URL Shortener</em></h1>
        <Tabs onSelect={index => setTask(indexMap.get(index)) }>
          <TabList>
            <Tab>Shorten URL</Tab>
            <Tab>Create URL</Tab>
            <Tab>Update URL</Tab>
          </TabList>
          <TabPanel>
            <ShortenerForm task={task}/>
          </TabPanel>
          <TabPanel>
            <ShortenerForm task={task}/>
          </TabPanel>
          <TabPanel>
            <ShortenerForm task={task}/>
          </TabPanel>
        </Tabs>
    </Grid>
  </ThemeProvider>
)};

export default App;
