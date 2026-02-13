package ui

const (
	maxContainerWidth = 100
	baseListWidth     = 35
	basePreviewWidth  = 50
	gapWidth          = 2
	paneBorderWidth   = 1
	containerPadding  = 2
	containerBorder   = 2
	paneHeight        = 20
)

type layoutMetrics struct {
	containerWidth int
	contentWidth   int
	listWidth      int
	previewWidth   int
}

func calculateLayout(termWidth int) layoutMetrics {
	if termWidth <= 0 {
		naturalContent := baseListWidth + basePreviewWidth + gapWidth + paneBorderWidth*2
		return layoutMetrics{
			containerWidth: naturalContent + containerPadding*2 + containerBorder,
			contentWidth:   naturalContent,
			listWidth:      baseListWidth,
			previewWidth:   basePreviewWidth,
		}
	}

	containerWidth := termWidth
	if containerWidth > maxContainerWidth {
		containerWidth = maxContainerWidth
	}

	innerWidth := containerWidth - (containerPadding*2 + containerBorder)
	if innerWidth < 4 {
		innerWidth = 4
	}

	usable := innerWidth - gapWidth - paneBorderWidth*2
	if usable < 2 {
		usable = 2
	}

	totalBase := baseListWidth + basePreviewWidth
	listWidth := usable * baseListWidth / totalBase
	previewWidth := usable - listWidth

	if listWidth < 1 {
		listWidth = 1
	}
	if previewWidth < 1 {
		previewWidth = 1
	}

	contentWidth := innerWidth

	return layoutMetrics{
		containerWidth: containerWidth,
		contentWidth:   contentWidth,
		listWidth:      listWidth,
		previewWidth:   previewWidth,
	}
}
