import React, {PropTypes} from 'react'
import TickscriptHeader from 'src/kapacitor/components/TickscriptHeader'
import FancyScrollbar from 'shared/components/FancyScrollbar'
import TickscriptEditor from 'src/kapacitor/components/TickscriptEditor'

const Tickscript = ({
  source,
  onSave,
  task,
  validation,
  onSelectDbrps,
  onChangeScript,
  onChangeType,
  isEditingID,
  onStartEditID,
  onStopEditID,
}) => (
  <div className="page">
    <TickscriptHeader
      task={task}
      source={source}
      onSave={onSave}
      isEditing={isEditingID}
      onStopEdit={onStopEditID}
      onStartEdit={onStartEditID}
      onChangeType={onChangeType}
      onSelectDbrps={onSelectDbrps}
    />
    <FancyScrollbar className="page-contents fancy-scroll--kapacitor">
      <div className="container-fluid">
        <div className="row">
          <div className="col-xs-12">
            {validation}
          </div>
        </div>
        <div className="row">
          <div className="col-xs-12">
            <TickscriptEditor
              script={task.script}
              onChangeScript={onChangeScript}
            />
          </div>
        </div>
      </div>
    </FancyScrollbar>
  </div>
)

const {arrayOf, bool, func, shape, string} = PropTypes

Tickscript.propTypes = {
  onSave: func.isRequired,
  source: shape(),
  task: shape({
    id: string,
    script: string,
    dbsrps: arrayOf(shape()),
  }).isRequired,
  onChangeScript: func.isRequired,
  onSelectDbrps: func.isRequired,
  validation: string,
  onChangeType: func.isRequired,
  isEditingID: bool.isRequired,
  onStartEditID: func.isRequired,
  onStopEditID: func.isRequired,
}

export default Tickscript
