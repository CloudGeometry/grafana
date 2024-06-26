import { Controller, useFormContext, useWatch } from 'react-hook-form';

import { DataSourceInstanceSettings } from '@grafana/data';
import { Field, FieldSet } from '@grafana/ui';
import { Trans, t } from 'app/core/internationalization';
import { DataSourcePicker } from 'app/features/datasources/components/picker/DataSourcePicker';

import { QueryEditorField } from './QueryEditorField';
import { useCorrelationsFormContext } from './correlationsFormContext';
import { FormDTO } from './types';

export const ConfigureCorrelationTargetForm = () => {
  const { control, formState } = useFormContext<FormDTO>();
  const withDsUID = (fn: Function) => (ds: DataSourceInstanceSettings) => fn(ds.uid);
  const { correlation } = useCorrelationsFormContext();
  const targetUID: string | undefined = useWatch({ name: 'targetUID' }) || correlation?.targetUID;

  return (
    <>
      <FieldSet label={t('correlations.target-form.title', 'Setup the target for the correlation (Step 2 of 3)')}>
        <Trans i18nKey="correlations.target-form.sub-text">
          <p>
            Define what data source the correlation will link to, and what query will run when the correlation is
            clicked.
          </p>
        </Trans>
        <Controller
          control={control}
          name="targetUID"
          rules={{
            required: { value: true, message: t('correlations.target-form.control-rules', 'This field is required.') },
          }}
          render={({ field: { onChange, value } }) => (
            <Field
              label={t('correlations.target-form.target-label', 'Target')}
              description={t(
                'correlations.target-form.target-description',
                'Specify which data source is queried when the link is clicked'
              )}
              htmlFor="target"
              invalid={!!formState.errors.targetUID}
              error={formState.errors.targetUID?.message}
            >
              <DataSourcePicker
                onChange={withDsUID(onChange)}
                noDefault
                current={value}
                inputId="target"
                width={32}
                disabled={correlation !== undefined}
              />
            </Field>
          )}
        />

        <QueryEditorField
          name="config.target"
          dsUid={targetUID}
          invalid={!!formState.errors?.config?.target}
          error={formState.errors?.config?.target?.message}
        />
      </FieldSet>
    </>
  );
};
