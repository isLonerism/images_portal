import React from 'react';
import './App.css';
import UploadTar from './UploadTar';
import HandleTar from './HandleTar';
import HandleImages from './HandleImages';
import PushImages from './PushImages';
import { makeStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import StepContent from '@material-ui/core/StepContent';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Box from '@material-ui/core/Box';
import AppBar from '@material-ui/core/AppBar';
import { Container, Toolbar, Tooltip } from '@material-ui/core';
import initENVNamespace from './environment.js';

class App extends React.Component {
  constructor () {
    super();
    this.state = {
      // Initialize activeStep variables with the first step
      activeStep: 0,
      // Updates when an error occures
      stepError: false,
      // Tar file name in object storage
      s3key: null,
      // Pod token used to identify pod between image load and push
      podToken: '',
      // List of images loaded from the tar file
      tarImageList: [],
      // List of the images received from the tar combined with 
      // the new chosen tag for each image
      taggedImageList: []
    };
    
    // Step Styling
    this.classes = makeStyles(theme => ({

      root: {
        width: '90%',
      },
      button: {
        marginTop: theme.spacing(1),
        marginRight: theme.spacing(1),
      },
      actionsContainer: {
        marginBottom: theme.spacing(2),
      },
      resetContainer: {
        padding: theme.spacing(3),
      },
    }));

    // Initialize ENV namespace from environment variables
    initENVNamespace();

    // Initialize steps
    this.steps = this.getSteps();

    // Bind functions
    this.handleTarUpload = this.handleTarUpload.bind(this);
    this.retrievePodToken = this.retrievePodToken.bind(this);
    this.retrieveTaggedImageList = this.retrieveTaggedImageList.bind(this);
    this.retrieveTarImageList = this.retrieveTarImageList.bind(this);
    this.handleNext = this.handleNext.bind(this);
    this.failStep = this.failStep.bind(this);
  };

  // Move on to the next step
  handleNext() {
    this.setState({
      activeStep: this.state.activeStep + 1
    });
  }

  // Handle reset - not impleneted
  handleReset() {
    console.log('reset')
  }

  // Returns all step titles
  getSteps() {
    return [
      <Tooltip
        title='Upload a tar archive containing one or 
               more images (usually created with "docker save" command) to the object storage'
        placement="top">
        <Typography variant="button">Upload Tar to Object Storage</Typography>
      </Tooltip>,
      <Tooltip
        title='Load the images from the tar archive to docker'
        placement="top">
        <Typography variant="button">Load Images to Docker</Typography>
      </Tooltip>,
      <Tooltip
        title="Tag the images loaded from the tar archive. 
               Specify once your OpenShift project's name,
               And then modify the name and version of each image (if needed) according to your
               needs, in the following format: 'name:version'.
               Review your images before clicking on the checkbox - cannot undo changes!"
        placement="top">
        <Typography variant="button">Tag Images</Typography>
      </Tooltip>,
      <Tooltip
        title="Push the tagged images into the OpenShift Registry.
               After a successful push, an Image Stream will be created for each image (can be
               viewed in the console under Builds -> Images)."
        placement="top">
        <Typography variant="button">Push Images to the OpenShift Registry</Typography>
      </Tooltip>
    ];
  }

  // Returns given step content
  getStepContent(step) {
    switch (step) {
      case 0:
        return  <Box width="100%" textAlign="center" pt="5%" pb="2%">
                  <UploadTar onTarUpload={this.handleTarUpload} failStep={this.failStep}></UploadTar>
                </Box>;
      case 1:
        return  <Box width="100%" textAlign="center" pt="5%" pb="2%">
                  <HandleTar s3key={this.state.s3key} podToken={this.retrievePodToken} tarImageList={this.retrieveTarImageList} failStep={this.failStep}></HandleTar>
                  </Box>;
      case 2:
        return  <Box width="100%" textAlign="center" pt="5%" pb="2%">
                  <HandleImages imageList={this.state.tarImageList} taggedImageList={this.retrieveTaggedImageList} failStep={this.failStep}></HandleImages>
                  </Box>;
      case 3:
        return  <Box width="100%" textAlign="center" pt="5%" pb="2%">
                  <PushImages podToken={this.state.podToken} imageList={this.state.taggedImageList} failStep={this.failStep} pushSuccess={this.handleNext}></PushImages>
                </Box>;
      default:
        return 'Unknown step';
    }
  }

  failStep(errDesc) {
    this.setState({ stepError: true});

    console.log(errDesc);
  };

  // Updates when a tar file is uploaded to s3 object storage, and 
  handleTarUpload = (uploadedTarFile) => {
    this.setState({ s3key: uploadedTarFile }, this.handleNext);
  };

  // Updates when a tar file is uploaded to s3 object storage
  retrievePodToken = (podToken) => {
    this.setState({ podToken: podToken });
  };

  // Updates when the user tags all images loaded from the tar file
  retrieveTaggedImageList = (imageList) => {
    this.setState({ taggedImageList: imageList }, this.handleNext());
  };

  // Initialize the list of images loaded from the tar file, and continute to the next step!
  retrieveTarImageList = (imageList) => {
    this.setState({ tarImageList: imageList }, this.handleNext);
  };

  render () {
    return (
      <Container maxWidth="sm">
        <AppBar color="primary">
          <Toolbar style={{display: "grid"}}>
            <Typography variant="h6" align="center">
            PaaS Team Portal
            </Typography>
          </Toolbar>
        </AppBar>
        <Box marginTop="25%" border={1} borderColor="primary.main">
          <Stepper activeStep={this.state.activeStep} orientation="vertical">
            {this.steps.map((label, index) => (
              <Step key={label}>
                <StepLabel error={this.state.stepError === true && this.state.activeStep === index}>{label}</StepLabel>
                <StepContent>
                  {this.getStepContent(index)}
{/*                    <div className={this.classes.actionsContainer}>
                    <div>
                      <Button
                        variant="contained"
                        color="primary"
                        onClick={this.handleNext}
                        className={this.classes.button}
                      >
                      {this.state.activeStep === this.steps.length - 1 ? 'Finish' : 'Next'}
                      </Button>
                    </div>
                  </div> */}
                </StepContent>
              </Step>
            ))}
          </Stepper>
          {this.state.activeStep === this.steps.length && (
          <Paper square elevation={0} className={this.classes.resetContainer} align='center'>
            <Typography variant="button" color="primary" style={{fontWeight: "bold", fontSize: 20 }}>
              Images Pushed Successfully!
              <br></br><br></br>
            </Typography>
          </Paper>
          )}
        </Box>
      </Container>
    )
  }
}

export default App;