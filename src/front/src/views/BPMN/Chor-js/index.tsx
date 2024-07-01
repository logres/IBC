import React, { useEffect, useState, useRef } from 'react';
import ChoreoModeler from './chor-js/Modeler'; // Adjust the import based on your package structure
// import ChoreoModeler from 'chor-js/lib/Modeler'; // Adjust the import based on your package structure
import PropertiesPanelModule from 'bpmn-js-properties-panel';
import PropertiesProviderModule from './lib-provider/properties-provider'; // Adjust the import based on your package structure
import '../../../../src/assets/styles/app.less';  // 确保您已经配置了less-loader来处理less文件
import blankXml from '../../../../src/assets/bpmns/newDiagram.bpmn'; // Adjust the import based on your package structure
import Reporter from './lib-provider/validator/Validator.js';
import axios from 'axios';
import MainPage from './pop-up/MainPage.js'
import TestPaletteProvider from './lib-provider/external-elements'

const ChorJs = () => {

  const isModelerHandling = useRef(false);
  // const [modeler, setModeler] = useState(null);
  const modeler = useRef(null);
  const reporter = useRef(null);
  let isDirty = false;
  let lastFile = null;
  let isValidating = false;

  async function renderModel(newXml) {
    await modeler.current.importXML(newXml);
    isDirty = false;
  }
  const generatePanelListener = () => {
    const panels = Array.prototype.slice.call(
      document.getElementById('panel-toggle').children
    );
    const panelsToListener = {}
    panels.forEach(p => {
      panelsToListener[p] = () => {
        panels.forEach(otherPanel => {
          if (p === otherPanel && !p.classList.contains('active')) {
            p.classList.add('active');
            document.getElementById(p.dataset.togglePanel).classList.remove('hidden');
          } else {
            otherPanel.classList.remove('active');
            document.getElementById(otherPanel.dataset.togglePanel).classList.add('hidden');
          }
        });
      }
    });
    return [panels, panelsToListener]
  }

  // returns the file name of the diagram currently being displayed
  function diagramName() {
    if (lastFile) {
      return lastFile.name;
    }
    return 'diagram.bpmn';
  }


  const js_download_listener = async (e: MouseEvent): Promise<void> => {
    console.log('downloadLink clicked');
    const downloadLink = document.getElementById('js-download-diagram');
    const result = await modeler.current.saveXML({ format: true });
    downloadLink['href'] = 'data:application/bpmn20-xml;charset=UTF-8,' + encodeURIComponent(result.xml);
    downloadLink['download'] = diagramName();
    isDirty = false;
  };

  const js_download_svg_listerner = async (e: MouseEvent): Promise<void> => {
    console.log('downloadSvgLink clicked');
    const downloadSvgLink = document.getElementById('js-download-svg');
    const result = await modeler.current.saveSVG();
    downloadSvgLink['href'] = 'data:image/svg+xml;charset=UTF-8,' + encodeURIComponent(result.svg);
    downloadSvgLink['download'] = diagramName() + '.svg';
  };

  const js_open_file_listener = (e: MouseEvent): void => {
    console.log('js-open-file clicked');
    document.getElementById('file-input').click();
  };

  const js_file_input_listener = (e: Event): void => {
    console.log('file-input changed');
    const loadDiagram = document.getElementById('file-input');
    const file = loadDiagram.files[0];
    if (file) {
      const reader = new FileReader();
      lastFile = file;
      reader.addEventListener('load', async () => {
        await renderModel(reader.result);
        loadDiagram.value = null; // allows reloading the same file
      }, false);
      reader.readAsText(file);
    }
  };

  const js_create_new_diagram_listener = async (e: MouseEvent): Promise<void> => {
    console.log('newDiagram clicked');
    await renderModel(blankXml);
    lastFile = false;
  };

  const drag_over_listener = (e: DragEvent): void => {
    const dropZone = document.body;
    e.preventDefault();
    dropZone.classList.add('is-dragover');
  };

  const drag_leave_listener = (e: DragEvent): void => {
    const dropZone = document.body;
    e.preventDefault();
    dropZone.classList.remove('is-dragover');
  };

  const is_drag_over_listener = (e: DragEvent): void => {
    const dropZone = document.body;
    e.preventDefault();
    dropZone.classList.remove('is-dragover');
    const file = e.dataTransfer.files[0];
    if (file) {
      const reader = new FileReader();
      lastFile = file;
      reader.addEventListener('load', () => {
        renderModel(reader.result);
      }, false);
      reader.readAsText(file);
    }
  };



  const js_validate_listener = (e: MouseEvent): void => {
    console.log('validateButton clicked');
    const validateButton = document.getElementById('js-validate');
    isValidating = !isValidating;
    if (isValidating) {
      reporter.current.validateDiagram();
      validateButton.classList.add('selected');
      validateButton['title'] = 'Disable checking';
    } else {
      reporter.current.clearAll();
      validateButton.classList.remove('selected');
      validateButton['title'] = 'Check diagram for problems';
    }
  };



  const js_upload_listener = async (e: MouseEvent): Promise<void> => {
    const bpmnName = prompt("请输入BPMN文件的名字：");
    if (bpmnName) {
      const confirmUpload = confirm("是否上传该bpmn文件？");
      if (confirmUpload) {
        const result = await modeler.current.saveXML({ format: true });
        console.log(result);

        const resultOfSvg = await modeler.current.saveSVG();
        console.log(resultOfSvg);

        var params = {};
        var queryString = window.location.search.substring(1);
        var pairs = queryString.split("&");
        for (var i = 0; i < pairs.length; i++) {
          var pair = pairs[i].split("=");
          params[pair[0]] = decodeURIComponent(pair[1]);
        }

        console.log("consortiumid = " + params["consortiumid"]);
        // console.log("userid = " + params["userid"])
        upload_bpmn_post(result, params, bpmnName, resultOfSvg);
      } else {
      }
    }
  };


  function upload_bpmn_post(result, params, bpmnName, resultOfSvg) {
    return axios.post('http://192.168.1.177:9999/chaincode/getPartByBpmnC', {
      bpmnContent: result.xml
    })
      .then((response) => {
        console.log('Post getParticipant request success:', response.data);
        axios.post(`http://192.168.1.177:8000/api/v1/consortiums/${params["consortiumid"]}/bpmns/_upload`, {
          bpmnContent: result.xml,
          consortiumid: params["consortiumid"],
          orgid: params["orgid"],
          name: bpmnName + '.bpmn', //之后加一个输入名字？
          svgContent: resultOfSvg.svg,
          participants: response.data
        })
          .then(function (response) {
            console.log(response);
          })
          .catch(function (error) {
            console.error(error);
          });
      })
      .catch((error) => {
        console.error('Post request error:', error);
      });
  }


  useEffect(() => {
    //download diagram as BPMN 
    const downloadLink = document.getElementById('js-download-diagram');
    console.log('downloadLink listener added');
    downloadLink.addEventListener('click', js_download_listener);


    // download diagram as SVG
    const downloadSvgLink = document.getElementById('js-download-svg');
    downloadSvgLink.addEventListener('click', js_download_svg_listerner);

    // open file dialog
    console.log('js_open_file add event listener');
    document.getElementById('js-open-file').addEventListener('click', js_open_file_listener);


    // load diagram from disk
    console.log('js_file_input add event listener');
    const loadDiagram = document.getElementById('file-input');
    loadDiagram.addEventListener('change', js_file_input_listener);

    // create new diagram
    const newDiagram = document.getElementById('js-new-diagram');
    newDiagram.addEventListener('click', js_create_new_diagram_listener);

    // drag & drop file
    const dropZone = document.body;
    dropZone.addEventListener('dragover', drag_over_listener);
    dropZone.addEventListener('dragleave', drag_leave_listener);
    dropZone.addEventListener('drop', is_drag_over_listener);

    // validation logic and toggle
    const validateButton = document.getElementById('js-validate');
    validateButton.addEventListener('click', js_validate_listener);

    //upload bpmn file
    document.getElementById('js-upload').addEventListener('click', js_upload_listener);

    window.addEventListener('beforeunload', function (e) {
      if (isDirty) {
        // see https://developer.mozilla.org/en-US/docs/Web/API/WindowEventHandlers/onbeforeunload
        e.preventDefault();
        e.returnValue = '';
      }
    });

    //modeler listener
    if (modeler.current === null || modeler.current === undefined) {
      return;
    }
    modeler.current.on('commandStack.changed', () => {
      if (isValidating) {
        reporter.current.validateDiagram();
      }
      isDirty = true;
    });
    modeler.current.on('import.render.complete', () => {
      if (isValidating) {
        reporter.current.validateDiagram();
      }
    });

    return () => {
      document.getElementById('js-download-diagram').removeEventListener('click', js_download_listener);
      document.getElementById('js-download-svg').removeEventListener('click', js_download_svg_listerner);
      document.getElementById('js-open-file').removeEventListener('click', js_open_file_listener);
      document.getElementById('file-input').removeEventListener('change', js_file_input_listener);
      document.getElementById('js-new-diagram').removeEventListener('click', js_create_new_diagram_listener);
      document.getElementById('js-validate').removeEventListener('click', js_validate_listener);
      document.getElementById('js-upload').removeEventListener('click', js_upload_listener);
      dropZone.removeEventListener('dragover', drag_over_listener);
      dropZone.removeEventListener('dragleave', drag_leave_listener);
      dropZone.removeEventListener('drop', is_drag_over_listener);
      modeler.current.off('commandStack.changed', on_change);
      console.log('event listeners[js-download-diagram, js-download-svg, js-open-file, file-input, js-new-diagram, js-validate, js-upload] removed');
    };
  }, []);

  useEffect(() => {

    // let modeler = null;
    const [panels, panelListeners] = generatePanelListener();

    const InitModeler = async () => {
      while (isModelerHandling.current) {
        await new Promise(resolve => setTimeout(resolve, 10));
      }
      console.log("Start InitModeler");
      isModelerHandling.current = true;
      modeler.current = new ChoreoModeler({
        container: '#canvas',
        propertiesPanel: {
          parent: '#properties-panel'
        },
        additionalModules: [
          PropertiesPanelModule,
          PropertiesProviderModule,
          TestPaletteProvider
        ],
        keyboard: {
          bindTo: document
        }
      })

      console.log('Init modeler', modeler.current);
      reporter.current = new Reporter(modeler.current);

      // toggle side panels

      console.log('panels', panels)
      panels.forEach(
        panel => panel.addEventListener('click', panelListeners[panel])
      )

      const initialRender = async () => {
        await renderModel(blankXml);
      }
      await initialRender();
      console.log("modeler initialized");
      window.bpmnjs = modeler.current;
      isModelerHandling.current = false
    }

    InitModeler();

    // You can add additional setup or model loading here
    return () => {
      const removeModeler = async () => {
        while (isModelerHandling.current) {
          await new Promise(resolve => setTimeout(resolve, 10));
        }
        isModelerHandling.current = true;
        console.log("Start removeModeler");
        console.log('Remove modeler', modeler.current)
        if (modeler.current !== null && modeler.current !== undefined) {
          await modeler.current.destroy()
          console.log('modeler destroyed')
        }
        // remove event listeners
        panels.forEach(
          panel => panel.removeEventListener('click', panelListeners[panel])
        )
        isModelerHandling.current = false;
      }
      removeModeler();
    };
  }, []);

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "calc(100vh - 160px)",
        background: "white",
      }}
    >
      <div className="content">
        <div id="canvas" style={{ height: '100%', width: '100%' }}></div>
        <div id="panel-toggle">
          <div data-toggle-panel="properties-panel" title="Toggle properties panel"><span>Properties</span></div>
        </div>
        <div id="properties-panel" className="side-panel hidden" ></div>
        <div className="buttons djs-container">
          <button id="js-new-diagram" className="icon-doc-new" title="Create new empty diagram"></button>
          <button id="js-open-file" className="icon-folder" title="Select BPMN XML file"></button>
          <div className="divider"></div>
          <a id="js-download-diagram" className="icon-file-code" title="Download BPMN XML file"></a>
          <a id="js-download-svg" className="icon-file-image" title="Download as SVG image"></a>
          <div className="divider"></div>
          <button id="js-validate" className="icon-bug" title="Check diagram for problems"></button>
          <div className="divider"></div>
          <button id="js-upload" className="icon-file-upload" title="Upload BPMN file"></button>
          <input id="file-input" name="name" type="file" accept=".bpmn, .xml" style={{ display: "none" }} />
        </div>
      </div>
      <MainPage />
    </div >
  );
};
export default ChorJs;